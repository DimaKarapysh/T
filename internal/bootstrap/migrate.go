package bootstrap

import (
	"T/internal/client"
	"T/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Migrate(lc fx.Lifecycle, config *config.Config, logger *zap.Logger) {
	err := client.RunMigrations(lc, config, logger)
	if err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}
}
