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

type config struct {
	Log         bool    `mapstructure:"log"`
	Debug       bool    `mapstructure:"debug"`
	Feature     Feature `mapstructure:"feature"`
	ListenAddr  string  `mapstructure:"addr"`
	BehindProxy bool    `mapstructure:"behind_proxy"`
	ExternalURL string  `mapstructure:"external_url"`
}

var _ configuration.Configuration = (*config)(nil)

func (c *config) Register(flagSet *pflag.FlagSet) {
	flagSet.String("http.external_url", "http://127.0.0.1:8080", "external url")
	flagSet.Bool("http.log", false, "enable http log")
	flagSet.Bool("http.debug", false, "enable http debug")
	flagSet.String("http.addr", ":8080", "http listen address")
	flagSet.Uint16("http.feature", uint16(FeatureAll), "http feature")
	flagSet.Bool("http.behind_proxy", false, "http behind proxy")
	common.MustNoError(viper.BindPFlags(flagSet))
	configuration.Register(c)
}

func (c *config) Read() {
	common.MustDecodeFromMapstructure(viper.AllSettings()["http"], c)
}
