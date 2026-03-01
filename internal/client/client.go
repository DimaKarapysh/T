package client

import (
	"T/internal/config"
	"T/migrations"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewPostgresClient(
	lc fx.Lifecycle,
	cfg *config.Config,
	logger *zap.Logger,
) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	parseConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("pgx parse config failed", zap.Error(err))
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), parseConfig)
	if err != nil {
		logger.Error("failed to connect to PostgreSQL", zap.Error(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing postgres pool")
			pool.Close()
			return nil
		},
	})

	logger.Info("connected to PostgreSQL")
	return pool, nil
}

func RunMigrations(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB for migration: %w", err)
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	logger.Info("мigrации успешно применены")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return nil
}

func NewLogger(logLevel string, serviceName string) *zap.Logger {
	level := zapcore.InfoLevel

	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "errwrap":
		level = zapcore.ErrorLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "message",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	logWriter := zapcore.Lock(os.Stdout)

	core := zapcore.NewCore(consoleEncoder, logWriter, level)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger.With(zap.String("service_name", serviceName))
}

func NewRedisClient(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	opts := &redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Error("failed to connect to Redis", zap.Error(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing Redis client")
			return client.Close()
		},
	})

	logger.Info("connected to Redis")
	return client, nil
}
