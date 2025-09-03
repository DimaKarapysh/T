package adapters

import (
	"T/internal/entity"
	"context"

	"github.com/google/uuid"
)

type SubService interface {
	Create(ctx context.Context, sub *entity.Task) (id uuid.UUID, err error)
	Get(ctx context.Context) (tasks []*entity.Task, err error)
	GetByID(ctx context.Context, id uuid.UUID) (task *entity.Task, err error)
	Update(ctx context.Context, id uuid.UUID, sub *entity.Task) error
	Delete(ctx context.Context, id uuid.UUID) (err error)
	List(ctx context.Context, limit, offset int) (tasks []*entity.Task, err error)
}
