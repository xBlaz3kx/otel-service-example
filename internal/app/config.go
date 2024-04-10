package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/observabilityExample/internal/pkg/config"
)

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
