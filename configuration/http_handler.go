package configuration

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/handlers/http"
)

var _ Configuration = (*httpHandlerConfiguration)(nil)
var HttpHandlerConfiguration = &httpHandlerConfiguration{}

type httpHandlerConfiguration struct {
	http.Config
}

func (*httpHandlerConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.Bool("http.log", false, "enable http log")
	flagSet.Bool("http.debug", false, "enable http debug")
	flagSet.String("http.addr", ":8080", "http listen address")
	flagSet.Uint16("http.feature", uint16(http.FeatureAll), "http feature")
	if err := viper.BindPFlags(flagSet); err != nil {
		panic(err)
	}
	Register(HttpHandlerConfiguration)
}

func (c *httpHandlerConfiguration) Read() {
	err := common.DecodeFromMapstructure(viper.AllSettings()["http"], &c.Config)
	if err != nil {
		panic(err)
	}
}
