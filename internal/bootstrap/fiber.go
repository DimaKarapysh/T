package bootstrap

import (
	"T/internal/config"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiber(
	lc fx.Lifecycle,
	cfg *config.Config,
	logger *zap.Logger,

) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	// lifecycle hook
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := fmt.Sprintf(":%s", cfg.RestPort)
				if err := app.Listen(addr); err != nil {
					logger.Fatal("HTTP server failed", zap.Error(err))
				}
			}()
			logger.Info("HTTP server started", zap.String("port", cfg.RestPort))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server...")
			return app.Shutdown()
		},
	})

	return app
}
