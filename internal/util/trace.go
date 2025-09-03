package util

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Step struct {
	span   trace.Span
	name   string
	logger *zap.Logger
}

// Start запускает спан. Можно передать атрибуты сразу.
func Start(ctx context.Context, tr trace.Tracer, logger *zap.Logger, name string, attrs ...attribute.KeyValue) (context.Context, *Step) {
	if tr == nil {
		tr = otel.Tracer("default")
	}

	ctx, span := tr.Start(ctx, name, trace.WithAttributes(attrs...))

	if logger != nil {
		logger.Debug("span started", zap.String("name", name))
	}

	return ctx, &Step{span: span, logger: logger, name: name}
}

// End завершает спан. Если err != nil — проставит статус Error и запишет ошибку.
func (s *Step) End(err error) {
	if s == nil || s.span == nil {
		return
	}
	defer s.span.End()

	fields := []zap.Field{zap.String("span", s.name)}

	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
		s.span.SetAttributes(
			attribute.Bool("err", true),
			attribute.String("err.msg", err.Error()),
		)
		var ce interface{ Code() string }
		if errors.As(err, &ce) {
			s.span.SetAttributes(attribute.String("error.code", ce.Code()))
		}

		if s.logger != nil {
			s.logger.Error("span failed", append(fields, zap.Error(err))...)
		}
	} else {
		s.span.SetStatus(codes.Ok, "success")
		if s.logger != nil {
			s.logger.Debug("span succeeded", fields...)
		}
	}

}

// Fail помечает текущий спан ошибкой, но НЕ завершает его (удобно между операциями).
func (s *Step) Fail(err error) {
	if s == nil || s.span == nil || err == nil {
		return
	}
	s.span.RecordError(err)
	s.span.SetStatus(codes.Error, err.Error())
}

// Event добавляет событие внутрь спана (быстрые отметки по ходу).
func (s *Step) Event(name string, attrs ...attribute.KeyValue) {
	if s == nil || s.span == nil {
		return
	}

	s.span.AddEvent(name, trace.WithAttributes(attrs...))

	if s.logger != nil {
		s.logger.Debug("span event", zap.String("span", s.name), zap.String("event", name))
	}

}
