package postgres

import (
	"T/internal/entity"
	"T/internal/errwrap"
	"T/internal/repo/postgres/sqlc"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type RepoSub struct {
	queries *sqlc.Queries
	logg    *zap.Logger
}

func NewRepoSub(queries *sqlc.Queries, logg *zap.Logger) *RepoSub {
	return &RepoSub{
		queries: queries,
		logg:    logg,
	}
}

func (a *RepoSub) Create(ctx context.Context, sub *entity.Subscription) error {
	const op = "CreateSub"
	logger := a.logg.With(zap.String("op", op), zap.String("user_id", sub.UserID.String()))

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	logger.Info("создание подписки")

	_, err := a.queries.CreateSubscription(ctx, sqlc.CreateSubscriptionParams{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       int32(sub.Price),
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	})
	if err != nil {
		return errwrap.SQLExecError(op, err)
	}

	logger.Info("подписка создана")
	return nil
}

func (a *RepoSub) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	const op = "GetByID"
	logger := a.logg.With(zap.String("op", op), zap.String("id", id.String()))

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	if id == uuid.Nil {
		return nil, errwrap.ErrUserIDEmpty(op)
	}

	logger.Info("подписка получается по id")

	dbUser, err := a.queries.GetSubscriptionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errwrap.NotFound(op, errors.New("subscription_not_found"))
		}
		return nil, errwrap.SQLQueryError(op, err)
	}

	logger.Info("подписка получена")

	return &entity.Subscription{
		ID:          dbUser.ID,
		ServiceName: dbUser.ServiceName,
		Price:       int(dbUser.Price),
		UserID:      dbUser.UserID,
		StartDate:   dbUser.StartDate,
		EndDate:     dbUser.EndDate,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
	}, nil
}

func (a *RepoSub) Update(ctx context.Context, s *entity.Subscription) error {
	const op = "RepoSub.Update"
	logger := a.logg.With(zap.String("op", op), zap.String("id", s.ID.String()))

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	logger.Info("обновление подписки")

	_, err := a.queries.UpdateSubscription(ctx, sqlc.UpdateSubscriptionParams{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       int32(s.Price),
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
	})
	if err != nil {
		return errwrap.SQLExecError(op, err)
	}

	logger.Info("подписка обновлена")
	return nil
}

func (a *RepoSub) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "RepoSub.Delete"
	logger := a.logg.With(zap.String("op", op), zap.String("id", id.String()))

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	logger.Info("удаление подписки")

	err := a.queries.DeleteSubscription(ctx, id)
	if err != nil {
		return errwrap.SQLExecError(op, err)
	}

	logger.Info("подписка помечена как удалённая")
	return nil
}

func (a *RepoSub) List(ctx context.Context, limit, offset int) ([]*entity.Subscription, error) {
	const op = "RepoSub.List"
	logger := a.logg.With(zap.String("op", op))

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	logger.Info("получение списка подписок")

	dbSubs, err := a.queries.ListSubscriptions(ctx, sqlc.ListSubscriptionsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, errwrap.SQLQueryError(op, err)
	}

	var result []*entity.Subscription
	for _, s := range dbSubs {
		result = append(result, &entity.Subscription{
			ID:          s.ID,
			ServiceName: s.ServiceName,
			Price:       int(s.Price),
			UserID:      s.UserID,
			StartDate:   s.StartDate,
			EndDate:     s.EndDate,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	logger.Info("подписки получены")
	return result, nil
}

func (a *RepoSub) Get(ctx context.Context) ([]*entity.Subscription, error) {
	const op = "RepoSub.Get"
	logger := a.logg.With(zap.String("op", op))

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	logger.Info("получение списка подписок")

	dbSubs, err := a.queries.GetSubscriptions(ctx)
	if err != nil {
		return nil, errwrap.SQLQueryError(op, err)
	}

	var result []*entity.Subscription
	for _, s := range dbSubs {
		result = append(result, &entity.Subscription{
			ID:          s.ID,
			ServiceName: s.ServiceName,
			Price:       int(s.Price),
			UserID:      s.UserID,
			StartDate:   s.StartDate,
			EndDate:     s.EndDate,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	logger.Info("подписки получены")
	return result, nil
}

func (a *RepoSub) Total(ctx context.Context, filter *entity.TotalFilter) ([]*entity.Subscription, error) {
	const op = "RepoSub.Total"
	logger := a.logg.With(zap.String("op", op))

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	logger.Info("получения массива подписок")
	//костыли с таймстампом
	arr, err := a.queries.GetTotalSubscription(ctx, sqlc.GetTotalSubscriptionParams{
		StartDateFrom: filter.StartFrom,
		StartDateTo: func() pgtype.Timestamp {
			if filter.EndTo != nil {
				var ts pgtype.Timestamp
				_ = ts.Scan(*filter.EndTo)
				return ts
			}
			return pgtype.Timestamp{Valid: false}
		}(),
		ServiceName: filter.ServiceName,
		UserID:      filter.UserID,
	})
	if err != nil {
		return nil, errwrap.SQLQueryError(op, err)
	}

	logger.Info("подписки получены", zap.Any("total", arr))

	return func(db []sqlc.Subscription) []*entity.Subscription {
		result := make([]*entity.Subscription, len(arr))
		for i, s := range arr {
			result[i] = &entity.Subscription{
				ID:          s.ID,
				ServiceName: s.ServiceName,
				Price:       int(s.Price),
				UserID:      s.UserID,
				StartDate:   s.StartDate,
				EndDate:     s.EndDate,
				CreatedAt:   s.CreatedAt,
				UpdatedAt:   s.UpdatedAt,
			}
		}
		return result
	}(arr), nil
}

func (a *RepoSub) TotalSQL(ctx context.Context, filter *entity.TotalFilter) (int, error) {
	const op = "RepoSub.Total"
	logger := a.logg.With(zap.String("op", op))

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	logger.Info("подсчет общей стоимости подписок")
	//костыли с таймстампом
	total, err := a.queries.GetTotalSubscriptionCost(ctx, sqlc.GetTotalSubscriptionCostParams{
		StartDateFrom: filter.StartFrom,
		StartDateTo: func() pgtype.Timestamp {
			if filter.EndTo != nil {
				var ts pgtype.Timestamp
				_ = ts.Scan(*filter.EndTo)
				return ts
			}
			return pgtype.Timestamp{Valid: false}
		}(),
		ServiceName: filter.ServiceName,
		UserID:      filter.UserID,
	})
	if err != nil {
		return 0, errwrap.SQLQueryError(op, err)
	}

	logger.Info("стоимость подсчитана", zap.Int64("total", total))
	return int(total), nil
}
