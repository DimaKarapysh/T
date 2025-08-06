package middlewarepackage

import (
	"T/internal/token"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	sessionIDKey contextKey = "session_id"
)

type middleware struct {
	accessToken token.Token
	log         *zap.Logger
}

func NewJWTMiddleware(accessToken *token.AccessToken, log *zap.Logger) JWTMiddleware {
	return &middleware{
		accessToken: accessToken.Token,
		log:         log,
	}
}

func (m *middleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := getBearerToken(c)
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"errwrap": "missing access token",
			})
		}

		claims, err := m.accessToken.VerifyToken(tokenStr)
		if err != nil {
			m.log.Warn("invalid access token", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"errwrap": "invalid or expired access token",
			})
		}

		ctx := context.WithValue(c.Context(), userIDKey, claims.UserID.String())
		ctx = context.WithValue(ctx, sessionIDKey, claims.SessionID.String())
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func getBearerToken(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	const prefix = "Bearer "

	if authHeader == "" || !strings.HasPrefix(authHeader, prefix) {
		return ""
	}
	return strings.TrimSpace(authHeader[len(prefix):])
}

// --- Optional Getters ---
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

func GetSessionID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(sessionIDKey).(string)
	return id, ok
}
