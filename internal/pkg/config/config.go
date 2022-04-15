package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	// Config keys
	TracingAddress     = "tracing.address"
	LogAddress         = "logging.address"
	PrometheusAddress  = "prometheus.address"
	PrometheusEndpoint = "prometheus.endpoint"
	MetricsKey         = "metrics"
	DebugKey           = "debug"

	// Flag keys
	DebugFlag              = "debug"
	MetricsFlag            = "metrics"
	TracingAddressFlag     = "tracing-addr"
	LogAddressFlag         = "log-addr"
	PrometheusAddressFlag  = "prometheus-address"
	PrometheusEndpointFlag = "prom-endpoint"
)

// SetupEnv binds environment variables to the viper configuration.
// Env variable should begin with EXAMPLE prefix, e.g. EXAMPLE_DEBUG=true should set the debug mode on.
// Another example is EXAMPLE_TRACING_ADDRESS="localhost:1337".
func SetupEnv() {
	viper.SetEnvPrefix("example")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// SetDefaults set default values for each key
func SetDefaults() {
	viper.SetDefault(TracingAddress, "localhost:55680")
	viper.SetDefault(LogAddress, "localhost:514")
	viper.SetDefault(PrometheusAddress, "localhost:8080")
	viper.SetDefault(DebugKey, false)
	viper.SetDefault(MetricsKey, false)
}

// ReadConfig load configuration from file
func ReadConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	home, _ := os.Getwd()
	viper.AddConfigPath(home + "/config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
}
