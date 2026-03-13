package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

// REPRESENTS JWT TOKEN PAYLOAD
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// HANDLES ACCESS AND REFRESH TOKEN OPERATIONS
type TokenManager struct{}

// CONSTRUCTOR
func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

// GENERATE ACCESS TOKEN (SHORT-LIVED) FOR API AUTHORIZATION
func (tm *TokenManager) GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
	expirationStr := os.Getenv("ACCESS_TOKEN_EXPIRATION")

	if expirationStr == "" {
		expirationStr = "15m"
	}

	duration, err := time.ParseDuration(expirationStr)

	if err != nil {
		duration = 15 * time.Minute // DEFAULT FALLBACK
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_KEY")))
}

// GENERATE REFRESH TOKEN (LONG-LIVED) FOR TOKEN RENEWAL
func (tm *TokenManager) GenerateRefreshToken(userID uuid.UUID, email, role string) (string, error) {
	expirationStr := os.Getenv("REFRESH_TOKEN_EXPIRATION")

	if expirationStr == "" {
		expirationStr = "168h" // 7 DAYS
	}

	duration, err := time.ParseDuration(expirationStr)

	if err != nil {
		duration = 7 * 24 * time.Hour // DEFAULT FALLBACK
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_KEY")))
}

// VERIFY ACCESS TOKEN AND RETURN CLAIMS
func (tm *TokenManager) VerifyAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// ENSURE SIGNING METHOD IS HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("signature method is not valid")
		}

		return []byte(os.Getenv("ACCESS_TOKEN_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("access token is not valid")
}

// VERIFY REFRESH TOKEN AND RETURN CLAIMS
func (tm *TokenManager) VerifyRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// ENSURE SIGNING METHOD IS HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("signature method is not valid")
		}

		return []byte(os.Getenv("REFRESH_TOKEN_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("refresh token is not valid")
}
