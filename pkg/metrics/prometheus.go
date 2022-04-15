package metrics

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	ExampleMeter = "exampleMeter"
	JobCounter   = "job_counter"
)

var (
	cfg = prometheus.Config{}
)

type Prometheus struct {
	controller *controller.Controller
	cfg        prometheus.Config
	exporter   *prometheus.Exporter
	http       *gin.Engine

	reqCounter metric.Int64Counter
	jobCounter metric.Int64Counter
}

// NewPrometheus creates a new HTTP server with metrics endpoint
func NewPrometheus(controller *controller.Controller, config prometheus.Config, exporter *prometheus.Exporter, endpoint string) *Prometheus {
	// Configure metrics endpoint
	r := gin.New()
	r.GET(endpoint, prometheusHandler(exporter))

	// Create a meter
	meter := controller.Meter(ExampleMeter)

	// Create a counter (an instrument according to OpenTelemetry)
	jobCounter, err := meter.NewInt64Counter(JobCounter)
	if err != nil {
		return nil
	}

	return &Prometheus{
		controller: controller,
		cfg:        config,
		exporter:   exporter,
		http:       r,
		jobCounter: jobCounter,
	}
}

func (p *Prometheus) AddJobCompleted(ctx context.Context, traceId, spanId string) {
	// todo instead of this idiomatic way of increasing the counter, use channels to make it better and pluggable
	p.jobCounter.Add(ctx, 1, attribute.String("traceId", traceId), attribute.String("spanId", spanId))
}

// Start Run the HTTP server with metrics
func (p *Prometheus) Start(address string) {
	err := p.http.Run(address)
	if err != nil {
		log.WithError(err).Errorf("Error exposing prometheus")
	}
}

func prometheusHandler(exporter *prometheus.Exporter) gin.HandlerFunc {
	return func(c *gin.Context) {
		exporter.ServeHTTP(c.Writer, c.Request)
	}
}

// Setup exposes prometheus metrics.
func Setup(res *resource.Resource, endpoint string) *Prometheus {
	// Create a controller
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(cfg.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(res),
	)

	// Create an exporter
	exporter, err := prometheus.New(cfg, c)
	if err != nil {
		log.WithError(err).Fatalf("Failed to initialize prometheus exporter")
	}

	global.SetMeterProvider(exporter.MeterProvider())

	return NewPrometheus(c, cfg, exporter, endpoint)
}
