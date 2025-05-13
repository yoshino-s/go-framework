package http

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/utils"
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
	Feature     Feature `mapstructure:"feature"`
	ListenAddr  string  `mapstructure:"addr"`
	BehindProxy bool    `mapstructure:"behind_proxy"`
	ExternalURL string  `mapstructure:"external_url"`

	Otel            bool `mapstructure:"otel"`
	ResponseTraceId bool `mapstructure:"response_trace_id"`
}

var _ configuration.Configuration = (*config)(nil)

func (c *config) Register(flagSet *pflag.FlagSet) {
	flagSet.String("http.external_url", "http://127.0.0.1:8080", "external url")
	flagSet.Bool("http.log", false, "enable http log")
	flagSet.String("http.addr", ":8080", "http listen address")
	flagSet.Uint16("http.feature", uint16(FeatureAll), "http feature")
	flagSet.Bool("http.behind_proxy", false, "http behind proxy")
	flagSet.Bool("http.otel", false, "enable opentelemetry")
	flagSet.Bool("http.response_trace_id", false, "enable x-trace-id in response header")
	utils.MustNoError(viper.BindPFlags(flagSet))
	configuration.Register(c)
}

func (c *config) Read() {
	utils.MustDecodeFromMapstructure(viper.AllSettings()["http"], c)
}
