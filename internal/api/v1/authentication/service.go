package authentication

import (
	"klinik-pkp-api/internal/api/v1/user"
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"crypto/sha256"
	"encoding/hex"
)

// CURRENT INSTANCE
type Service struct {
	db           *gorm.DB
	validator    *utils.Validator
	tokenManager *utils.TokenManager
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
		db:           db,
		validator:    utils.NewValidator(),
		tokenManager: utils.NewTokenManager(),
	}
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *Service) findAuthentication(refreshToken string) (*Authentication, error) {
	var authentication Authentication

	hashedToken := hashRefreshToken(refreshToken)

	// NEW FORMAT: SHA-256 HASH
	if err := s.db.Where("token = ?", hashedToken).First(&authentication).Error; err == nil {
		return &authentication, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// LEGACY FALLBACK:
	// Allows sessions created before this security change to rotate once.
	if err := s.db.Where("token = ?", refreshToken).First(&authentication).Error; err != nil {
		return nil, err
	}

	return &authentication, nil
}

// VERIFIES CREDENTIALS, GENERATES BOTH TOKENS, STORES REFRESH TOKEN IN DB
func (s *Service) AddAuthentication(payload AddAuthenticationPayload) (fiber.Map, error) {
	var user user.User

	// FIND USER BY EMAIL
	if err := s.db.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		return nil, err
	}

	if payload.NIP != user.NIP {
		return nil, gorm.ErrRecordNotFound
	}

	// CHECK IF USER IS ACTIVE
	if !user.IsActive {
		return nil, gorm.ErrRecordNotFound
	}

	// VERIFY PASSWORD
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return nil, err
	}

	// GENERATE ACCESS TOKEN (SHORT-LIVED)
	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// GENERATE REFRESH TOKEN (LONG-LIVED)
	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(&Authentication{Token: hashRefreshToken(refreshToken)}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return fiber.Map{
		"id":            user.ID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

// VERIFIES REFRESH TOKEN, THEN GENERATES NEW ACCESS TOKEN
func (s *Service) UpdateAuthentication(refreshToken string) (fiber.Map, error) {
	// VERIFY SIGNATURE AND EXPIRATION FIRST
	claims, err := s.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// CHECK WHETHER TOKEN IS STILL AN ACTIVE SESSION
	authentication, err := s.findAuthentication(refreshToken)
	if err != nil {
		return nil, err
	}

	// GENERATE NEW ACCESS TOKEN
	accessToken, err := s.tokenManager.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return nil, err
	}

	// ROTATE REFRESH TOKEN
	newRefreshToken, err := s.tokenManager.GenerateRefreshToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return nil, err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// REVOKE OLD TOKEN ATOMICALLY.
	// RowsAffected prevents the same refresh token being reused concurrently.
	result := tx.Where("token = ?", authentication.Token).
		Delete(&Authentication{})

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if result.RowsAffected != 1 {
		tx.Rollback()
		return nil, gorm.ErrRecordNotFound
	}

	// STORE ONLY HASH OF NEW REFRESH TOKEN
	if err := tx.Create(&Authentication{
		Token: hashRefreshToken(newRefreshToken),
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return fiber.Map{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}, nil
}

// VERIFIES REFRESH TOKEN EXISTS IN DB, THEN REMOVES IT
func (s *Service) DeleteAuthentication(refreshToken string) error {
	hashedToken := hashRefreshToken(refreshToken)

	// Support both new hashed sessions and legacy plaintext sessions.
	result := s.db.
		Where("token = ? OR token = ?", hashedToken, refreshToken).
		Delete(&Authentication{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
