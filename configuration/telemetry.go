package configuration

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/telemetry"
)

var _ Configuration = (*telemetryConfiguration)(nil)
var TelemetryConfiguration = &telemetryConfiguration{}

type telemetryConfiguration struct {
	telemetry.Config
}

func (*telemetryConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.String("telemetry.sentry_dsn", "", "sentry dsn")
	if err := viper.BindPFlags(flagSet); err != nil {
		panic(err)
	}
	Register(TelemetryConfiguration)
}

func (c *telemetryConfiguration) Read() {
	err := common.DecodeFromMapstructure(viper.AllSettings()["telemetry"], &c.Config)
	if err != nil {
		panic(err)
	}
}
