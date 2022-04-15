package worker

import (
	"context"
	"fmt"
	"github.com/xBlaz3kx/observabilityExample/pkg/logger"
	"github.com/xBlaz3kx/observabilityExample/pkg/metrics"
	"github.com/xBlaz3kx/observabilityExample/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

var (
	commonLabels = []attribute.KeyValue{
		attribute.String("labelA", "test1"),
		attribute.String("labelB", "test2"),
		attribute.String("labelC", "test3"),
	}
)

type (
	Worker struct {
		tracer     *tracer.Tracer
		prometheus *metrics.Prometheus
		logger     *logger.Logger
	}
)

func NewWorker(tracer *tracer.Tracer, prometheus *metrics.Prometheus, logger *logger.Logger) *Worker {
	return &Worker{
		tracer:     tracer,
		prometheus: prometheus,
		logger:     logger,
	}
}

func (w *Worker) Start(ctx context.Context) {
	// Work begins; start a trace
	ctx, span := w.tracer.Start(
		ctx,
		"Parent",
		"Parent span!",
		trace.WithAttributes(commonLabels...))

	// Simulate a workload
	for i := 0; i < 10; i++ {
		// Start a new span, log both the traceId and spanId
		spanName := fmt.Sprintf("Sample-%d", i)
		_, iSpan := w.tracer.Start(ctx, spanName, "Doing really hard work")
		// Wait for a second
		<-time.After(time.Second)

		// Add an example event
		iSpan.AddEvent("example-event")

		if w.prometheus != nil {
			// Increase the number of jobs completed
			// To create an Exemplar, add a label with traceId to the metric
			w.prometheus.AddJobCompleted(ctx, iSpan.SpanContext().TraceID().String(), iSpan.SpanContext().SpanID().String())
		}

		iSpan.End()
	}

	span.End()

	w.logger.Get().Info("Done!")
}
