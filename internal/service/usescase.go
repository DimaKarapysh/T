package service

import (
	"T/internal/config"
	"T/internal/entity"
	"T/internal/errwrap"
	"T/internal/repo/ports"
	"T/internal/service/adapters"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type SubService struct {
	logger     *zap.Logger
	repository ports.Repository
	//sessionRepo ports.Session
	//token       token.Token
	config *config.Config
}

type SubParams struct {
	Logger     *zap.Logger
	Repository ports.Repository
	//SessionRepo ports.Session
	//Token       token.Token
	Config *config.Config
}

func NewSubService(subParams SubParams) adapters.SubService {
	return &SubService{logger: subParams.Logger, repository: subParams.Repository, config: subParams.Config}
}

func (s *SubService) Create(ctx context.Context, sub *entity.Subscription) (uuid.UUID, error) {
	const op = "SubService.Create"
	logger := s.logger.With(zap.String("op", op),
		zap.String("user_id", sub.UserID.String()))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Info("создание подписки", zap.Any("sub", sub))

	err := s.repository.Create(ctx, sub)
	return sub.ID, err
}

func (s *SubService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	const op = "SubService.GetByID"
	logger := s.logger.With(zap.String("op", op),
		zap.String("id", id.String()))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Info("получение подписки по ID")
	return s.repository.GetByID(ctx, id)
}

func (s *SubService) Update(ctx context.Context, id uuid.UUID, sub *entity.Subscription) error {
	const op = "SubService.Update"
	logger := s.logger.With(zap.String("op", op),
		zap.String("id", id.String()),
		zap.String("user_id", sub.UserID.String()))

	//ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	//defer cancel()

	sub.ID = id

	logger.Info("обновление подписки", zap.Any("sub", sub))

	err := s.repository.Update(ctx, sub)
	if err != nil {
		return errwrap.SQLQueryError(op, err)
	}
	return nil
}

func (s *SubService) Get(ctx context.Context) ([]*entity.Subscription, error) {
	const op = "SubService.List"
	logger := s.logger.With(zap.String("op", op))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Info("получение списка подписок")
	return s.repository.Get(ctx)
}

func (s *SubService) List(ctx context.Context, limit, offset int) ([]*entity.Subscription, error) {

	const op = "SubService.List"
	logger := s.logger.With(zap.String("op", op))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Info("получение списка подписок")
	list, err := s.repository.List(ctx, limit, offset)
	if err != nil {
		return nil, errwrap.SQLQueryError(op, err)
	}
	return list, nil
}

func (s *SubService) Total(ctx context.Context, filter *entity.TotalFilter) (int, error) {
	const op = "SubService.Total"
	logger := s.logger.With(zap.String("op", op), zap.String("user_id", filter.UserID.String()))

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	logger.Info("получение общей стоимости подписок")

	subs, err := s.repository.Total(ctx, filter)
	if err != nil {
		return 0, errwrap.SQLQueryError(op, err)
	}

	var total int
	for _, sub := range subs {
		total += sub.Price
	}

	logger.Info("подсчёт завершён", zap.Int("total", total))
	return total, nil

	//можно посчитать общую стоимость через SQL функцию
	//totalSQL, err := s.repository.TotalSQL(ctx, *filter)
	//if err != nil {
	//	return 0, errwrap.SQLQueryError(op, err)
	//}
	//logger.Info("подсчёт завершён", zap.Int("total", totalSQL))
	//return totalSQL, nil
}

func (s *SubService) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "SubService.Delete"
	logger := s.logger.With(zap.String("op", op),
		zap.String("id", id.String()))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Info("удаление подписки")

	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errwrap.SQLQueryError(op, err)
	}
	return nil
}
