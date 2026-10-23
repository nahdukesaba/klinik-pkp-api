package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type localAccessTokenVerifier interface {
	VerifyAccessToken(string) (*Claims, error)
}

type klinikSSOContextResolver interface {
	ResolveKlinikContext(context.Context, string) (*SSOContext, error)
}

type authMiddleware struct {
	local         localAccessTokenVerifier
	central       klinikSSOContextResolver
	authorization KlinikAuthorizationResolver
}

// AuthMiddleware accepts a legacy Klinik JWT or a Central SSO token, then
// resolves the current Klinik membership before trusting any role.
func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return newAuthMiddleware(
		NewTokenManager(),
		NewSSOClient(),
		NewKlinikAuthorizationResolver(db),
		false,
	)
}

// OptionalAuthMiddleware authenticates a supplied bearer token but allows a
// request without one. It is used by logout to derive trusted login context.
func OptionalAuthMiddleware(db *gorm.DB) fiber.Handler {
	return newAuthMiddleware(
		NewTokenManager(),
		NewSSOClient(),
		NewKlinikAuthorizationResolver(db),
		true,
	)
}

func newAuthMiddleware(
	local localAccessTokenVerifier,
	central klinikSSOContextResolver,
	authorization KlinikAuthorizationResolver,
	optional bool,
) fiber.Handler {
	middleware := &authMiddleware{
		local:         local,
		central:       central,
		authorization: authorization,
	}

	return func(ctx *fiber.Ctx) error {
		authHeader := strings.TrimSpace(ctx.Get("Authorization"))

		if authHeader == "" {
			if optional {
				return ctx.Next()
			}
			return authenticationError(ctx, fiber.StatusUnauthorized)
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {
			return authenticationError(ctx, fiber.StatusUnauthorized)
		}

		trusted, err := middleware.authenticate(ctx.UserContext(), parts[1])
		if err != nil {
			switch {
			case errors.Is(err, ErrSSOUnavailable), errors.Is(err, ErrKlinikAuthorizationUnavailable):
				return JSONResponse(ctx, fiber.StatusServiceUnavailable, "authentication service unavailable", fiber.ErrServiceUnavailable, true)
			case errors.Is(err, ErrKlinikAuthorizationDenied):
				return authenticationError(ctx, fiber.StatusForbidden)
			default:
				return authenticationError(ctx, fiber.StatusUnauthorized)
			}
		}

		setTrustedAuthContext(ctx, trusted)
		return ctx.Next()
	}
}

func (m *authMiddleware) authenticate(ctx context.Context, accessToken string) (*AuthContext, error) {
	claims, localErr := m.local.VerifyAccessToken(accessToken)
	if localErr == nil {
		trusted, err := m.authorization.ResolveLegacyKlinikAuthorization(ctx, claims.UserID)
		if err != nil {
			return nil, err
		}
		trusted.LoginSource = LoginSourceDirect
		return trusted, nil
	}

	ssoContext, err := m.central.ResolveKlinikContext(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if ssoContext.AppUserID == uuid.Nil ||
		(ssoContext.ID != uuid.Nil && ssoContext.ID != ssoContext.AppUserID) ||
		ssoContext.Status != "approved" ||
		!IsKlinikRole(strings.ToLower(strings.TrimSpace(ssoContext.Role))) {
		return nil, ErrSSOUnauthorized
	}

	trusted, err := m.authorization.ResolveKlinikAuthorization(ctx, ssoContext.AppUserID)
	if err != nil {
		return nil, err
	}
	if trusted.Role != strings.ToLower(strings.TrimSpace(ssoContext.Role)) {
		return nil, ErrKlinikAuthorizationDenied
	}
	trusted.LoginSource = loginSourceFromVerifiedJWT(accessToken)
	return trusted, nil
}

func setTrustedAuthContext(ctx *fiber.Ctx, trusted *AuthContext) {
	ctx.Locals("auth_context", trusted)
	contextID := trusted.AppUserID
	if trusted.LocalUserID != nil {
		contextID = *trusted.LocalUserID
	}
	ctx.Locals("id", contextID)
	ctx.Locals("app_user_id", trusted.AppUserID)
	ctx.Locals("email", trusted.Email)
	ctx.Locals("role", trusted.Role)
	ctx.Locals("login_source", string(trusted.LoginSource))
}

func authenticationError(ctx *fiber.Ctx, status int) error {
	fiberError := fiber.ErrUnauthorized
	if status == fiber.StatusForbidden {
		fiberError = fiber.ErrForbidden
	}
	return JSONResponse(ctx, status, fiberError.Message, fiberError, true)
}

func loginSourceFromVerifiedJWT(token string) LoginSource {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return LoginSourceDirect
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return LoginSourceDirect
	}
	var claims struct {
		AMR []json.RawMessage `json:"amr"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return LoginSourceDirect
	}

	for _, entry := range claims.AMR {
		var method string
		if err := json.Unmarshal(entry, &method); err != nil {
			var structured struct {
				Method string `json:"method"`
			}
			if json.Unmarshal(entry, &structured) != nil {
				continue
			}
			method = structured.Method
		}
		switch strings.ToLower(strings.TrimSpace(method)) {
		case "magiclink", "otp", "sso", "saml":
			return LoginSourceSSO
		}
	}
	return LoginSourceDirect
}

// ROLE MIDDLEWARE - RESTRICTS ACCESS TO SPECIFIC ROLES.
// MUST BE USED AFTER AUTH MIDDLEWARE.
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		role, ok := ctx.Locals("role").(string)
		if !ok {
			return JSONResponse(
				ctx,
				fiber.StatusForbidden,
				fiber.ErrForbidden.Message,
				fiber.ErrForbidden,
				true,
			)
		}

		if slices.Contains(allowedRoles, role) {
			return ctx.Next()
		}

		return JSONResponse(
			ctx,
			fiber.StatusForbidden,
			fiber.ErrForbidden.Message,
			fiber.ErrForbidden,
			true,
		)
	}
}
