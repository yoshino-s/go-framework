package telemetry

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
)

type config struct {
	SentryDSN string `mapstructure:"sentry_dsn"`
}

var _ configuration.Configuration = (*telemetryConfiguration)(nil)

type telemetryConfiguration struct {
	config *config
}

func (t *telemetryConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.String("telemetry.sentry_dsn", "", "sentry dsn")
	if err := viper.BindPFlags(flagSet); err != nil {
		panic(err)
	}
	configuration.Register(t)
}

func (c *telemetryConfiguration) Read() {
	err := common.DecodeFromMapstructure(viper.AllSettings()["telemetry"], c.config)
	if err != nil {
		panic(err)
	}

	if c.config.SentryDSN != "" {
		initSentry(c.config.SentryDSN)
	}
}
