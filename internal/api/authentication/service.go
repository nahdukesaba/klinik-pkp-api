package authentication

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/user"
	"klinik-pkp-api/utils"
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

// PAYLOAD FOR REFRESHING ACCESS TOKEN
type RefreshAuthenticationPayload struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// PAYLOAD FOR DELETING AUTHENTICATION (LOGOUT)
type DeleteAuthenticationPayload struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:           db,
		validator:    utils.NewValidator(),
		tokenManager: utils.NewTokenManager(),
	}
}

// VERIFIES CREDENTIALS, GENERATES BOTH TOKENS, STORES REFRESH TOKEN IN DB
func (s *Service) AddRefreshToken(payload AddAuthenticationPayload) (*AddAuthenticationResponse, error) {
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

	// STORE REFRESH TOKEN IN DATABASE
	if err := s.db.Create(&Authentication{Token: refreshToken}).Error; err != nil {
		return nil, err
	}

	return &AddAuthenticationResponse{
		ID:           user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// VERIFIES REFRESH TOKEN, THEN GENERATES NEW ACCESS TOKEN
func (s *Service) UpdateAccessToken(payload RefreshAuthenticationPayload) (*RefreshAuthenticationResponse, error) {
	// CHECK IF REFRESH TOKEN EXISTS
	if err := s.db.Where("token = ?", payload.RefreshToken).First(&Authentication{}).Error; err != nil {
		return nil, err
	}

	// VERIFY REFRESH TOKEN SIGNATURE AND EXTRACT CLAIMS
	claims, err := s.tokenManager.VerifyRefreshToken(payload.RefreshToken)
	if err != nil {
		return nil, err
	}

	// GENERATE NEW ACCESS TOKEN
	accessToken, err := s.tokenManager.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return nil, err
	}

	return &RefreshAuthenticationResponse{
		AccessToken: accessToken,
	}, nil
}

// VERIFIES REFRESH TOKEN EXISTS IN DB, THEN REMOVES IT
func (s *Service) DeleteRefreshToken(payload DeleteAuthenticationPayload) error {
	result := s.db.Where("token = ?", payload.RefreshToken).Delete(&Authentication{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
