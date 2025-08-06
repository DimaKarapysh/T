package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	SessionID uuid.UUID `json:"session_id"`
	jwt.RegisteredClaims
}

func newClaims(userId uuid.UUID, roleID uuid.UUID, expirationTime time.Duration, sessionID uuid.UUID) *Claims {
	claims := &Claims{
		UserID:    userId,
		SessionID: sessionID,
		RoleID:    roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return claims
}
