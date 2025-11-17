package authn_oidc

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/utils"
)

type OidcAuthenticationConfig struct {
	IssuerURL    string            `mapstructure:"issuer_url"`
	ClientID     string            `mapstructure:"client_id"`
	ClientSecret string            `mapstructure:"client_secret"`
	RedirectURL  string            `mapstructure:"redirect_url"`
	Scopes       []string          `mapstructure:"scopes"`
	DefaultRole  string            `mapstructure:"default_role"`
	GroupRoleMap map[string]string `mapstructure:"group_role_map"`
}

var _ configuration.Configuration = (*OidcAuthenticationConfig)(nil)

func (c *OidcAuthenticationConfig) Register(flagSet *pflag.FlagSet) {
	flagSet.String("authn.oidc.issuer_url", "", "OIDC issuer URL")
	flagSet.String("authn.oidc.client_id", "", "OIDC client ID")
	flagSet.String("authn.oidc.client_secret", "", "OIDC client secret")
	flagSet.String("authn.oidc.redirect_url", "", "OIDC redirect URL")
	flagSet.StringSlice("authn.oidc.scopes", []string{"openid", "profile", "email", "groups"}, "OIDC scopes")
	flagSet.String("authn.oidc.default_role", "", "Default role for OIDC users")
	flagSet.StringToString("authn.oidc.group_role_map", nil, "Map of OIDC groups to application roles")
	utils.MustNoError(viper.BindPFlags(flagSet))
	configuration.Register(c)
}

func (c *OidcAuthenticationConfig) Read() {
	authn, ok := viper.AllSettings()["authn"]
	if !ok {
		panic("authn.jwt configuration not found")
	}

	utils.MustDecodeFromMapstructure(authn.(map[string]any)["oidc"], c)

	if c.GroupRoleMap != nil {
		lowerMap := make(map[string]string)
		for k, v := range c.GroupRoleMap {
			lowerMap[strings.ToLower(k)] = v
		}
		c.GroupRoleMap = lowerMap
	}
}
