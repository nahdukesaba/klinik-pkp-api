package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"strings"
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

// SET REFRESH TOKEN COOKIE WITH SECURE ATTRIBUTES FOR BROWSER-BASED SESSION RENEWAL
func (tm *TokenManager) SetRefreshTokenCookie(ctx *fiber.Ctx, refreshToken string) {
	secure, sameSite := setRefreshTokenSecurity()
	expirationStr := os.Getenv("REFRESH_TOKEN_EXPIRATION")

	if expirationStr == "" {
		expirationStr = "168h"
	}

	expiration, err := time.ParseDuration(expirationStr)
	if err != nil {
		expiration = 7 * 24 * time.Hour // DEFAULT FALLBACK
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   int(expiration.Seconds()),
		Expires:  time.Now().Add(expiration),
	})
}

// EXPIRE REFRESH TOKEN COOKIE DURING LOGOUT
func (tm *TokenManager) ClearRefreshTokenCookie(ctx *fiber.Ctx) {
	secure, sameSite := setRefreshTokenSecurity()

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

// LOAD COOKIE SECURITY OPTIONS FROM ENV WITH SAFE FALLBACK
func setRefreshTokenSecurity() (bool, string) {
	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	secure := strings.ToLower(strings.TrimSpace(os.Getenv("REFRESH_COOKIE_SECURE"))) == "true" || appEnv == "production"
	sameSite := strings.ToLower(strings.TrimSpace(os.Getenv("REFRESH_COOKIE_SAMESITE")))

	if sameSite == "" {
		if secure {
			sameSite = fiber.CookieSameSiteNoneMode
		} else {
			sameSite = fiber.CookieSameSiteLaxMode
		}
	}

	switch sameSite {
	case fiber.CookieSameSiteNoneMode, fiber.CookieSameSiteLaxMode, fiber.CookieSameSiteStrictMode, fiber.CookieSameSiteDisabled:
	default:
		if secure {
			sameSite = fiber.CookieSameSiteNoneMode
		} else {
			sameSite = fiber.CookieSameSiteLaxMode
		}
	}

	// BROWSERS REJECT SameSite=None WHEN Secure IS FALSE.
	if sameSite == fiber.CookieSameSiteNoneMode && !secure {
		sameSite = fiber.CookieSameSiteLaxMode
	}

	return secure, sameSite
}
