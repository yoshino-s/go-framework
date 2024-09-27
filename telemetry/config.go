package telemetry

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
)

var _ configuration.Configuration = (*telemetryConfiguration)(nil)

type telemetryConfiguration struct {
	SentryDSN          string  `mapstructure:"sentry_dsn"`
	TracesSampleRate   float64 `mapstructure:"traces_sample_rate"`
	ProfilesSampleRate float64 `mapstructure:"profiles_sample_rate"`
}

func (t *telemetryConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.String("telemetry.sentry_dsn", "", "sentry dsn")
	flagSet.Float64("telemetry.traces_sample_rate", 1.0, "traces sample rate")
	flagSet.Float64("telemetry.profiles_sample_rate", 1.0, "profiles sample rate")
	common.MustNoError(viper.BindPFlags(flagSet))
	configuration.Register(t)
}

func (c *telemetryConfiguration) Read() {
	common.MustNoError(common.DecodeFromMapstructure(viper.AllSettings()["telemetry"], c))

	if c.SentryDSN != "" {
		c.initSentry()
	}
}
