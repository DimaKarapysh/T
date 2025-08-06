package bootstrap

import (
	"T/internal/client"
	"T/internal/config"
	"T/internal/repo/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newPostgres(lc fx.Lifecycle, config *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	return client.NewPostgresClient(lc,
		config,
		logger,
	)
}

func newSQLCQueries(pool *pgxpool.Pool) *sqlc.Queries {
	return sqlc.New(pool)
}
