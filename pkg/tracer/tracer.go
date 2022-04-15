package tracer

import (
	"context"
	"github.com/xBlaz3kx/observabilityExample/pkg/logger"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	logger *logger.Logger
	tracer trace.Tracer
}

func (t *Tracer) Start(ctx context.Context, spanName, log string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(
		ctx,
		spanName,
		opts...)

	spanCtx := span.SpanContext()
	t.logger.WithTrace(spanCtx.TraceID().String(), spanCtx.SpanID().String()).Info(log)
	return ctx, span
}

func NewTracer(logger *logger.Logger, tracer trace.Tracer) *Tracer {
	return &Tracer{
		logger: logger,
		tracer: tracer,
	}
}
