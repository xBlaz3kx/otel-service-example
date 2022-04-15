package app

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/observabilityExample/internal/pkg/config"
	"github.com/xBlaz3kx/observabilityExample/internal/service/worker"
	"github.com/xBlaz3kx/observabilityExample/pkg/logger"
	"github.com/xBlaz3kx/observabilityExample/pkg/metrics"
	"github.com/xBlaz3kx/observabilityExample/pkg/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"os"
	"os/signal"
)

const (
	ServiceName = "test-service"
	tracerName  = "test-tracer"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "example",
		Short: "",
		Long:  ``,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	var (
		// Logging
		loggingAddress = viper.GetString(config.LogAddress)
		// Tracing
		tracingAddress = viper.GetString(config.TracingAddress)
		// Metrics
		metricsAddress  = viper.GetString(config.LogAddress)
		metricsEndpoint = viper.GetString(config.LogAddress)
		enableMetrics   = viper.GetBool(config.MetricsKey)
		isDebug         = viper.GetBool(config.DebugKey)

		ctx         = context.Background()
		quitChannel = make(chan os.Signal)
		traceLogger = logger.NewLogger(loggingAddress)
		prometheus  *metrics.Prometheus

		// Configure global resource for observability
		res, err = resource.New(ctx,
			resource.WithAttributes(
				semconv.ServiceNameKey.String(ServiceName),
			),
			resource.WithOS(),
			resource.WithHost(),
		)
	)

	signal.Notify(quitChannel, os.Interrupt)

	if err != nil {
		log.WithError(err).Fatalf("Cannot create resource")
	}

	// If it is debug mode, set level to trace
	if isDebug {
		log.SetLevel(log.TraceLevel)
	}

	// Connect to the tracing backend
	shutdown := tracer.InitProvider(res, tracingAddress)
	defer shutdown()

	// Create a new tracer
	newTracer := tracer.NewTracer(traceLogger, otel.Tracer(tracerName))

	// If metrics are enabled, expose a prometheus endpoint
	if enableMetrics {
		prometheus = metrics.Setup(res, metricsEndpoint)
		go prometheus.Start(metricsAddress)
	}

	// Start a worker
	exampleWorker := worker.NewWorker(newTracer, prometheus, traceLogger)
	go exampleWorker.Start(ctx)

	<-quitChannel
}

func initConfig() {
	config.SetDefaults()
	config.SetupEnv()
	config.ReadConfig(cfgFile)
}

func setFlags() {
	// Base
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolP(config.DebugFlag, "d", false, "debug mode")
	// Metric configuration
	rootCmd.PersistentFlags().BoolP(config.MetricsFlag, "m", false, "enable metrics")
	rootCmd.PersistentFlags().String(config.PrometheusAddressFlag, "", "prometheus address")
	rootCmd.PersistentFlags().String(config.PrometheusEndpointFlag, "", "metrics endpoint")
	// Tracing
	rootCmd.PersistentFlags().String(config.TracingAddressFlag, "", "tracing address")
	// Logging
	rootCmd.PersistentFlags().String(config.LogAddressFlag, "", "log address")

	// Bind flags to viper configuration
	_ = viper.BindPFlag(config.PrometheusAddress, rootCmd.Flags().Lookup(config.PrometheusAddressFlag))
	_ = viper.BindPFlag(config.TracingAddress, rootCmd.Flags().Lookup(config.TracingAddressFlag))
	_ = viper.BindPFlag(config.LogAddress, rootCmd.Flags().Lookup(config.LogAddressFlag))
	_ = viper.BindPFlag(config.PrometheusEndpoint, rootCmd.Flags().Lookup(config.PrometheusEndpointFlag))
	_ = viper.BindPFlag(config.MetricsKey, rootCmd.Flags().Lookup(config.MetricsFlag))
	_ = viper.BindPFlag(config.DebugKey, rootCmd.Flags().Lookup(config.DebugFlag))
}

func init() {
	cobra.OnInitialize(initConfig)
}

func Execute() {
	setFlags()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error executing the command")
	}
}
