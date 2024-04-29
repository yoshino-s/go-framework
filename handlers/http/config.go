package http

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
)

type Feature uint16

const (
	FeatureNone    Feature = 0
	FeatureVersion Feature = 1 << iota
	FeatureHealth
	FeatureReady
	FeatureMetrics

	FeatureAll = FeatureVersion | FeatureHealth | FeatureReady | FeatureMetrics
)

func (f Feature) Has(flag Feature) bool {
	return f&flag != 0
}

func (f Feature) Add(flag Feature) Feature {
	return f | flag
}

func (f Feature) Remove(flag Feature) Feature {
	return f &^ flag
}

type Config struct {
	Log        bool    `mapstructure:"log"`
	Debug      bool    `mapstructure:"debug"`
	Feature    Feature `mapstructure:"feature"`
	ListenAddr string  `mapstructure:"addr"`
}

var DefaultConfig = Config{
	Log:        true,
	Debug:      false,
	Feature:    FeatureAll,
	ListenAddr: ":8080",
}

var _ configuration.Configuration = (*httpHandlerConfiguration)(nil)

type httpHandlerConfiguration struct {
	config *Config
}

func (c *httpHandlerConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.Bool("http.log", false, "enable http log")
	flagSet.Bool("http.debug", false, "enable http debug")
	flagSet.String("http.addr", ":8080", "http listen address")
	flagSet.Uint16("http.feature", uint16(FeatureAll), "http feature")
	common.MustNoError(viper.BindPFlags(flagSet))
	configuration.Register(c)
}

func (c *httpHandlerConfiguration) Read() {
	common.MustDecodeFromMapstructure(viper.AllSettings()["http"], c.config)
}
