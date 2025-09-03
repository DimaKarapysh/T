package service

import (
	"T/internal/config"
	"T/internal/entity"
	"T/internal/repo/ports"
	"T/internal/service/adapters"
	"T/internal/util"
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const Tracer = "usecase.task"

type SubService struct {
	logger     *zap.Logger
	repository ports.Repository
	config     *config.Config
	tracer     trace.Tracer
}

type SubParams struct {
	Logger     *zap.Logger
	Repository ports.Repository
	Config     *config.Config
}

func NewSubService(subParams SubParams) adapters.SubService {
	return &SubService{
		logger:     subParams.Logger,
		repository: subParams.Repository,
		config:     subParams.Config,
		tracer:     otel.Tracer(Tracer),
	}
}

func (s *SubService) Create(ctx context.Context, sub *entity.Task) (id uuid.UUID, err error) {
	const op = "SubService.Create"
	logger := s.logger.With(zap.String("op", op), zap.String("id", sub.ID.String()))

	ctx, step := util.Start(ctx, s.tracer, logger, "Service:IsCreated",
		attribute.String("id", sub.ID.String()),
	)
	defer step.End(err)

	logger.Info("создание задачи", zap.Any("task", sub))

	id, err = s.repository.Create(ctx, sub)
	if err != nil {
		logger.Error(op, zap.Error(err))
		return uuid.Nil, nil
	}
	return id, nil
}

func (s *SubService) GetByID(ctx context.Context, id uuid.UUID) (task *entity.Task, err error) {
	const op = "SubService.GetByID"
	logger := s.logger.With(zap.String("op", op), zap.String("id", id.String()))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ctx, step := util.Start(ctx, s.tracer, logger, "Service:GetByID",
		attribute.String("id", id.String()),
	)
	defer step.End(err)

	logger.Info("получение задачи по ID")

	task, err = s.repository.GetByID(ctx, id)
	if err != nil {
		logger.Error(op, zap.Error(err))
		return nil, err
	}

	logger.Info("задача получена")
	return task, nil
}

func (s *SubService) Update(ctx context.Context, id uuid.UUID, sub *entity.Task) (err error) {
	const op = "SubService.Update"
	logger := s.logger.With(
		zap.String("op", op),
		zap.String("id", id.String()),
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// нормализуем ID
	sub.ID = id

	ctx, step := util.Start(ctx, s.tracer, logger, "Service:IsUpdated",
		attribute.String("id", id.String()),
	)
	defer step.End(err)

	logger.Info("обновление задачи", zap.Any("task", sub))

	err = s.repository.Update(ctx, sub)
	if err != nil {
		logger.Error(op, zap.Error(err))
		return err
	}
	return nil
}

func (s *SubService) Get(ctx context.Context) (tasks []*entity.Task, err error) {
	const op = "SubService.Get"
	logger := s.logger.With(zap.String("op", op))

	ctx, step := util.Start(ctx, s.tracer, logger, "Service:GetAll")
	defer step.End(err)

	logger.Info("получение всех задач")

	tasks, err = s.repository.Get(ctx)
	if err != nil {
		logger.Error(op, zap.Error(err))
		return nil, err
	}
	return tasks, nil
}

func (s *SubService) List(ctx context.Context, limit, offset int) (tasks []*entity.Task, err error) {
	const op = "SubService.List"
	logger := s.logger.With(zap.String("op", op))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ctx, step := util.Start(ctx, s.tracer, logger, "Service:List",
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	)
	defer step.End(err)

	logger.Info("получение списка задач", zap.Int("limit", limit), zap.Int("offset", offset))

	tasks, err = s.repository.List(ctx, limit, offset)
	if err != nil {
		logger.Error(op, zap.Error(err))
		return nil, err
	}
	return tasks, nil
}

func (s *SubService) Delete(ctx context.Context, id uuid.UUID) (err error) {
	const op = "SubService.Delete"
	logger := s.logger.With(zap.String("op", op), zap.String("id", id.String()))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ctx, step := util.Start(ctx, s.tracer, logger, "Service:IsDeleted",
		attribute.String("id", id.String()),
	)
	defer step.End(err)

	logger.Info("удаление задачи")

	err = s.repository.Delete(ctx, id)
	if err != nil {
		logger.Error(op, zap.Error(err))
		return err
	}
	return nil
}
