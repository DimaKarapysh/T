package postgres

import (
	"T/internal/entity"
	"T/internal/errwrap"
	"T/internal/repo/postgres/sqlc"
	"T/internal/util"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type RepoSub struct {
	queries *sqlc.Queries
	logg    *zap.Logger
	trace   trace.Tracer
}

func NewRepoSub(queries *sqlc.Queries, logg *zap.Logger, tr trace.Tracer) *RepoSub {
	return &RepoSub{
		queries: queries,
		logg:    logg,
		trace:   tr,
	}
}

func (a *RepoSub) Create(ctx context.Context, sub *entity.Task) (id uuid.UUID, err error) {
	const op = "CreateSub"
	logger := a.logg.With(zap.String("op", op), zap.String("task_id", sub.ID.String()))

	ctx, step := util.Start(ctx, a.trace, logger, "Repo:Create")
	var spanErr error
	defer step.End(spanErr)

	logger.Info("создание задачи")

	task, err := a.queries.CreateTask(ctx, sqlc.CreateTaskParams{
		Title: sub.Title,
		Description: func(p *string) pgtype.Text {
			if p == nil {
				return pgtype.Text{Valid: false}
			}
			return pgtype.Text{String: *p, Valid: true}
		}(sub.Description),
		Status: pgtype.Text{String: string(sub.Status), Valid: true},
	})
	if err != nil {
		spanErr = err
		return uuid.Nil, errwrap.SQLExecError(op, err)
	}

	logger.Info("задача создана")
	return task.ID, nil
}

func (a *RepoSub) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	const op = "GetByID"
	logger := a.logg.With(zap.String("op", op), zap.String("id", id.String()))

	ctx, step := util.Start(ctx, a.trace, logger, "Repo:GetByID")
	var spanErr error
	defer step.End(spanErr)

	if id == uuid.Nil {
		err := errwrap.ErrUserIDEmpty(op)
		spanErr = err
		return nil, err
	}

	logger.Info("получение задачи по id")

	row, err := a.queries.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			spanErr = err
			return nil, errwrap.NotFound(op, errors.New("task_not_found"))
		}
		spanErr = err
		return nil, errwrap.SQLQueryError(op, err)
	}

	logger.Info("задача получена")

	task := &entity.Task{
		ID:          row.ID,
		Title:       row.Title,
		Description: textPtr(row.Description),
		Status:      entity.Status(row.Status.String),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
		DeleteAt:    timePtr(row.DeletedAt),
	}
	return task, nil
}

func (a *RepoSub) Update(ctx context.Context, s *entity.Task) error {
	const op = "RepoSub.Update"
	logger := a.logg.With(zap.String("op", op), zap.String("id", s.ID.String()))

	ctx, step := util.Start(ctx, a.trace, logger, "Repo:Update")
	var spanErr error
	defer step.End(spanErr)

	logger.Info("обновление задачи")

	_, err := a.queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:    s.ID,
		Title: s.Title,
		Description: func(p *string) pgtype.Text {
			if p == nil {
				return pgtype.Text{Valid: false}
			}
			return pgtype.Text{String: *p, Valid: true}
		}(s.Description),
		Status: pgtype.Text{String: string(s.Status), Valid: true},
	})
	if err != nil {
		spanErr = err
		return errwrap.SQLExecError(op, err)
	}

	logger.Info("задача обновлена")
	return nil
}

func (a *RepoSub) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "RepoSub.Delete"
	logger := a.logg.With(zap.String("op", op), zap.String("id", id.String()))

	ctx, step := util.Start(ctx, a.trace, logger, "Repo:Delete")
	var spanErr error
	defer step.End(spanErr)

	logger.Info("удаление задачи")

	if err := a.queries.DeleteTask(ctx, id); err != nil {
		spanErr = err
		return errwrap.SQLExecError(op, err)
	}

	logger.Info("задача помечена как удалённая")
	return nil
}

func (a *RepoSub) List(ctx context.Context, limit, offset int) ([]*entity.Task, error) {
	const op = "RepoSub.List"
	logger := a.logg.With(zap.String("op", op))

	ctx, step := util.Start(ctx, a.trace, logger, "Repo:List")
	var spanErr error
	defer step.End(spanErr)

	logger.Info("получение списка задач")

	rows, err := a.queries.ListTasks(ctx, sqlc.ListTasksParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		spanErr = err
		return nil, errwrap.SQLQueryError(op, err)
	}

	result := make([]*entity.Task, 0, len(rows))
	for _, r := range rows {
		result = append(result, &entity.Task{
			ID:          r.ID,
			Title:       r.Title,
			Description: textPtr(r.Description),
			Status:      entity.Status(r.Status.String),
			CreatedAt:   r.CreatedAt.Time,
			UpdatedAt:   r.UpdatedAt.Time,
		})
	}

	logger.Info("задачи получены")
	return result, nil
}

func (a *RepoSub) Get(ctx context.Context) ([]*entity.Task, error) {
	const op = "RepoSub.Get"
	logger := a.logg.With(zap.String("op", op))

	ctx, step := util.Start(ctx, a.trace, logger, "Repo:GetAll")
	var spanErr error
	defer step.End(spanErr)

	logger.Info("получение всех задач")

	rows, err := a.queries.GetTasks(ctx)
	if err != nil {
		spanErr = err
		return nil, errwrap.SQLQueryError(op, err)
	}

	result := make([]*entity.Task, 0, len(rows))
	for _, r := range rows {
		result = append(result, &entity.Task{
			ID:          r.ID,
			Title:       r.Title,
			Description: textPtr(r.Description),
			Status:      entity.Status(r.Status.String),
			CreatedAt:   r.CreatedAt.Time,
			UpdatedAt:   r.UpdatedAt.Time,
		})
	}

	logger.Info("задачи получены")
	return result, nil
}

func textPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	v := t.String
	return &v
}

func timePtr(ts pgtype.Timestamp) *time.Time {
	if !ts.Valid {
		return nil
	}
	v := ts.Time
	return &v
}
