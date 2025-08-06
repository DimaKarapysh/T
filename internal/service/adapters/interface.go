package adapters

import (
	"T/internal/entity"
	"context"
	"github.com/google/uuid"
)

type SubService interface {
	Create(ctx context.Context, sub *entity.Subscription) (uuid.UUID, error)
	Get(ctx context.Context) ([]*entity.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, sub *entity.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.Subscription, error)
	Total(ctx context.Context, filter *entity.TotalFilter) (int, error)
}
