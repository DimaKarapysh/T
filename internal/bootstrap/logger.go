package bootstrap

import (
	"T/internal/client"
	"T/internal/config"
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func InitLogger(lc fx.Lifecycle, config *config.Config) *zap.Logger {
	logger := client.NewLogger(config.LoggerLevel, config.AppName)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Sync()
			return nil
		},
	})
	return logger
}
