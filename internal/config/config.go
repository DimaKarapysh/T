package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type JaegerConfig struct {
	// без http:// и без /api/traces — просто host:port для gRPC
	Endpoint string `envconfig:"JAEGER_ENDPOINT" default:"localhost:4317"`
}
type JWTConfig struct {
	Secret     string        `envconfig:"JWT_SECRET" required:"true"`
	AccessTTL  time.Duration `envconfig:"APP_TOKEN_TIME_LIVE" default:"30m"`
	RefreshTTL time.Duration `envconfig:"APP_REFRESH_TIME_LIVE" default:"24h"`
}

type PostgresConfig struct {
	Host     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port     string `envconfig:"POSTGRES_PORT" default:"5432"`
	User     string `envconfig:"POSTGRES_USER" default:"postgres"`
	Password string `envconfig:"POSTGRES_PASSWORD" default:"password"`
	Database string `envconfig:"POSTGRES_DATABASE" default:"auth"`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	Database int    `envconfig:"REDIS_DB" default:"0"`
}

type Config struct {
	AppName     string `envconfig:"APP_NAME" default:"auth_service"`
	RestPort    string `envconfig:"REST_PORT" default:"8080"`
	LoggerLevel string `envconfig:"LOGGER_LEVEL" default:"debug"`
	WebhookURL  string `envconfig:"WEBHOOK_URL"`
	JWT         JWTConfig
	Postgres    PostgresConfig
	Redis       RedisConfig
	Jaeger      JaegerConfig
}

func Load() *Config {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("не удалось загрузить конфигурацию: %v", err)
	}

	return &cfg
}
