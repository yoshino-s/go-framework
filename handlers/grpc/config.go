package grpc

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/utils"
)

type config struct {
	ListenAddr string `mapstructure:"addr"`
}

var _ configuration.Configuration = (*config)(nil)

func (c *config) Register(flagSet *pflag.FlagSet) {
	flagSet.String("grpc.addr", "127.0.0.1:50051", "grpc listen address")
	configuration.Register(c)
}

func (c *config) Read() {
	utils.MustDecodeFromMapstructure(viper.AllSettings()["grpc"], c)
}
