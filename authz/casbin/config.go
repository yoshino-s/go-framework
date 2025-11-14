package authz_casbin

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/utils"
)

type CasbinAuthorizationConfig struct {
	ModelConfPath  string `mapstructure:"model_conf_path"`
	PolicyConfPath string `mapstructure:"policy_conf_path"`
}

var _ configuration.Configuration = (*CasbinAuthorizationConfig)(nil)

func (c *CasbinAuthorizationConfig) Register(flagSet *pflag.FlagSet) {
	flagSet.String("authz.casbin.model_conf_path", "", "Casbin model configuration file path")
	flagSet.String("authz.casbin.policy_conf_path", "", "Casbin policy configuration file path")
	utils.MustNoError(viper.BindPFlags(flagSet))
	configuration.Register(c)
}

func (c *CasbinAuthorizationConfig) Read() {
	authz, ok := viper.AllSettings()["authz"]
	if !ok {
		return
	}

	utils.MustDecodeFromMapstructure(authz.(map[string]any)["casbin"], c)
}
