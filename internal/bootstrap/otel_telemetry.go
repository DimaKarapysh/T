package bootstrap

import (
	"T/internal/config"
	"context"
	"time"

	"go.opentelemetry.io/otel"
	jaegerExporter "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newOpenTelemetry(lc fx.Lifecycle, conf *config.Config, logg *zap.Logger) (*sdktrace.TracerProvider, error) {
	exp, err := jaegerExporter.New(jaegerExporter.WithCollectorEndpoint(
		jaegerExporter.WithEndpoint(conf.Jaeger.Endpoint),
	))
	if err != nil {
		logg.Error("jaeger exporter init failed", zap.Error(err))
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp,
			sdktrace.WithMaxExportBatchSize(10),
			sdktrace.WithBatchTimeout(2*time.Second),
		),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(conf.AppName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return tp.Shutdown(ctx)
		},
	})

	return tp, nil
}

func newTracer(conf *config.Config) trace.Tracer {
	return otel.Tracer(conf.AppName)
}
