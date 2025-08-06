package bootstrap

import (
	"T/internal/client"
	"T/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func redisInit(lc fx.Lifecycle, config *config.Config, logger *zap.Logger) (*redis.Client, error) {
	return client.NewRedisClient(lc,
		config,
		logger,
	)
}
