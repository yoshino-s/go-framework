package mock

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
)

type config struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

var _ configuration.Configuration = (*config)(nil)

func (c *config) Register(set *pflag.FlagSet) {
	set.String("mock_oidc.client_id", "", "Client ID")
	set.String("mock_oidc.client_secret", "", "Client Secret")
	common.MustNoError(viper.BindPFlags(set))
	configuration.Register(c)
}

func (c *config) Read() {
	common.MustDecodeFromMapstructure(viper.AllSettings()["mock_oidc"], c)
}
