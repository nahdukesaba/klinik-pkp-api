package authentication

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"klinik-pkp-api/internal/api/v1/user"
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Service struct {
	db            *gorm.DB
	store         legacyAuthenticationStore
	validator     *utils.Validator
	tokenManager  *utils.TokenManager
	authorization utils.KlinikAuthorizationResolver
}

type legacyAuthenticationStore interface {
	FindUserByEmail(context.Context, string) (*user.User, error)
	FindAuthentication(context.Context, string, string) (*Authentication, error)
	CreateAuthentication(context.Context, string) error
	RotateAuthentication(context.Context, string, string) error
	DeleteAuthentication(context.Context, string, string) error
}

type gormLegacyAuthenticationStore struct {
	db *gorm.DB
}

// PAYLOAD FOR LOGIN REQUEST
type AddAuthenticationPayload struct {
	Email    string `json:"email" validate:"required,email"`
	NIP      string `json:"nip" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:            db,
		store:         &gormLegacyAuthenticationStore{db: db},
		validator:     utils.NewValidator(),
		tokenManager:  utils.NewTokenManager(),
		authorization: utils.NewKlinikAuthorizationResolver(db),
	}
}

func newServiceWithDependencies(
	db *gorm.DB,
	store legacyAuthenticationStore,
	authorization utils.KlinikAuthorizationResolver,
) *Service {
	service := NewService(db)
	service.store = store
	service.authorization = authorization
	return service
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *gormLegacyAuthenticationStore) FindUserByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	var localUser user.User
	if s.db == nil {
		return nil, gorm.ErrRecordNotFound
	}
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&localUser).Error; err != nil {
		return nil, err
	}
	return &localUser, nil
}

func (s *gormLegacyAuthenticationStore) FindAuthentication(
	ctx context.Context,
	hashedToken string,
	plainToken string,
) (*Authentication, error) {
	var authentication Authentication
	if s.db == nil {
		return nil, gorm.ErrRecordNotFound
	}

	// NEW FORMAT: SHA-256 HASH
	if err := s.db.WithContext(ctx).Where("token = ?", hashedToken).First(&authentication).Error; err == nil {
		return &authentication, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// LEGACY FALLBACK:
	// Allows sessions created before this security change to rotate once.
	if err := s.db.WithContext(ctx).Where("token = ?", plainToken).First(&authentication).Error; err != nil {
		return nil, err
	}

	return &authentication, nil
}

func (s *gormLegacyAuthenticationStore) CreateAuthentication(ctx context.Context, token string) error {
	if s.db == nil {
		return gorm.ErrInvalidDB
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&Authentication{Token: token}).Error
	})
}

func (s *gormLegacyAuthenticationStore) RotateAuthentication(
	ctx context.Context,
	oldToken string,
	newToken string,
) error {
	if s.db == nil {
		return gorm.ErrInvalidDB
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("token = ?", oldToken).Delete(&Authentication{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		return tx.Create(&Authentication{Token: newToken}).Error
	})
}

func (s *gormLegacyAuthenticationStore) DeleteAuthentication(
	ctx context.Context,
	hashedToken string,
	plainToken string,
) error {
	if s.db == nil {
		return gorm.ErrInvalidDB
	}
	result := s.db.WithContext(ctx).
		Where("token = ? OR token = ?", hashedToken, plainToken).
		Delete(&Authentication{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// VERIFIES CREDENTIALS, GENERATES BOTH TOKENS, STORES REFRESH TOKEN IN DB
func (s *Service) AddAuthentication(ctx context.Context, payload AddAuthenticationPayload) (fiber.Map, error) {
	localUser, err := s.store.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	trusted, err := s.authorizeLegacyCredentials(ctx, payload, localUser)
	if err != nil {
		return nil, err
	}
	if trusted.LocalUserID == nil {
		return nil, utils.ErrKlinikAuthorizationDenied
	}

	// GENERATE ACCESS TOKEN (SHORT-LIVED)
	accessToken, err := s.tokenManager.GenerateAccessToken(
		*trusted.LocalUserID,
		trusted.Email,
		trusted.Role,
	)
	if err != nil {
		return nil, err
	}

	// GENERATE REFRESH TOKEN (LONG-LIVED)
	refreshToken, err := s.tokenManager.GenerateRefreshToken(
		*trusted.LocalUserID,
		trusted.Email,
		trusted.Role,
	)
	if err != nil {
		return nil, err
	}

	if err := s.store.CreateAuthentication(ctx, hashRefreshToken(refreshToken)); err != nil {
		return nil, err
	}

	return fiber.Map{
		"id":            *trusted.LocalUserID,
		"app_user_id":   trusted.AppUserID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func (s *Service) authorizeLegacyCredentials(
	ctx context.Context,
	payload AddAuthenticationPayload,
	localUser *user.User,
) (*utils.AuthContext, error) {
	if localUser == nil || payload.NIP != localUser.NIP || !localUser.IsActive {
		return nil, gorm.ErrRecordNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(localUser.Password), []byte(payload.Password)); err != nil {
		return nil, err
	}
	return s.authorization.ResolveLegacyKlinikAuthorization(ctx, localUser.ID)
}

// VERIFIES REFRESH TOKEN, THEN GENERATES NEW ACCESS TOKEN
func (s *Service) UpdateAuthentication(ctx context.Context, refreshToken string) (fiber.Map, error) {
	// VERIFY SIGNATURE AND EXPIRATION FIRST
	claims, err := s.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// CHECK WHETHER TOKEN IS STILL AN ACTIVE SESSION
	authentication, err := s.store.FindAuthentication(
		ctx,
		hashRefreshToken(refreshToken),
		refreshToken,
	)
	if err != nil {
		return nil, err
	}
	trusted, err := s.authorization.ResolveLegacyKlinikAuthorization(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if trusted.LocalUserID == nil {
		return nil, utils.ErrKlinikAuthorizationDenied
	}

	// GENERATE NEW ACCESS TOKEN
	accessToken, err := s.tokenManager.GenerateAccessToken(
		*trusted.LocalUserID,
		trusted.Email,
		trusted.Role,
	)
	if err != nil {
		return nil, err
	}

	// ROTATE REFRESH TOKEN
	newRefreshToken, err := s.tokenManager.GenerateRefreshToken(
		*trusted.LocalUserID,
		trusted.Email,
		trusted.Role,
	)
	if err != nil {
		return nil, err
	}

	if err := s.store.RotateAuthentication(
		ctx,
		authentication.Token,
		hashRefreshToken(newRefreshToken),
	); err != nil {
		return nil, err
	}

	return fiber.Map{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}, nil
}

// VERIFIES REFRESH TOKEN EXISTS IN DB, THEN REMOVES IT
func (s *Service) DeleteAuthentication(ctx context.Context, refreshToken string) error {
	return s.store.DeleteAuthentication(ctx, hashRefreshToken(refreshToken), refreshToken)
}
