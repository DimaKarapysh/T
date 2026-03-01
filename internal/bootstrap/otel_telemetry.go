package bootstrap

import (
	"context"
	"time"

	"T/internal/config"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func newOpenTelemetry(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) (*sdktrace.TracerProvider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Экспортер в Jaeger через OTLP gRPC
	exp, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(cfg.Jaeger.Endpoint), // jaeger:4317 / localhost:4317
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create otlp exporter")
	}

	// Ресурс с именем сервиса
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.AppName), // task_service
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create resource")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down tracer provider")
			return tp.Shutdown(ctx)
		},
	})

	return tp, nil
}

func newTracer(tp *sdktrace.TracerProvider) trace.Tracer {
	return tp.Tracer("task_service")
}
