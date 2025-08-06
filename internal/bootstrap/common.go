package bootstrap

import (
	"T/internal/config"
	"go.uber.org/fx"
)

func CommonModules() fx.Option {
	return fx.Options(
		fx.Provide(
			config.Load,

			newPostgres,
			newSQLCQueries,
			//redisInit,
			InitLogger,
			InitValidator,
			NewFiber,
		),
	)
}
