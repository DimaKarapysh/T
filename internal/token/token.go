package token

import (
	"github.com/google/uuid"
	"time"
)

type Token interface {
	CreateToken(userID, roleID, sessionID uuid.UUID, duration time.Duration) (string, error)
	VerifyToken(token string) (*Claims, error)
}
