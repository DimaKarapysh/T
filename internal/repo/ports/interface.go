package ports

import (
	"T/internal/entity"
	"T/internal/repo/redis"
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, s *entity.Task) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	Update(ctx context.Context, s *entity.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.Task, error)
	Get(ctx context.Context) ([]*entity.Task, error)
}

type Session interface {
	SetRefreshToken(ctx context.Context, userID uuid.UUID, token string, sessionID uuid.UUID, userAgent, ip string) (err error)
	GetRefreshToken(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (Session *redis.SessionData, err error)
	RevokeRefreshToken(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (err error)
	RevokeAllRefreshTokens(ctx context.Context, userID uuid.UUID) (err error)
}
