package oidc

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
)

type providerConfig struct {
	IssuerURL     string   `mapstructure:"issuer_url"`
	AuthURL       string   `mapstructure:"auth_url"`
	TokenURL      string   `mapstructure:"token_url"`
	DeviceAuthURL string   `mapstructure:"device_auth_url"`
	UserInfoURL   string   `mapstructure:"user_info_url"`
	JWKSURL       string   `mapstructure:"jwks_url"`
	Algorithms    []string `mapstructure:"algorithms"`
}

type config struct {
	ClientID       string         `mapstructure:"client_id"`
	ClientSecret   string         `mapstructure:"client_secret"`
	Scopes         []string       `mapstructure:"scopes"`
	ProviderConfig providerConfig `mapstructure:"provider_config"`
	IssuerURL      string         `mapstructure:"issuer_url"`
}

var _ configuration.Configuration = (*config)(nil)

func (c *config) Register(flagSet *pflag.FlagSet) {
	flagSet.String("oauth2.client_id", "", "oauth2 client id")
	flagSet.String("oauth2.client_secret", "", "oauth2 client secret")
	flagSet.StringSlice("oauth2.scopes", []string{oidc.ScopeOpenID, "profile", "email"}, "oauth2 scopes")
	flagSet.String("oauth2.issuer_url", "", "oauth2 issuer url")
	flagSet.String("oauth2.provider_config.issuer_url", "", "oauth2 provider issuer url")
	flagSet.String("oauth2.provider_config.auth_url", "", "oauth2 provider auth url")
	flagSet.String("oauth2.provider_config.token_url", "", "oauth2 provider token url")
	flagSet.String("oauth2.provider_config.device_auth_url", "", "oauth2 provider device auth url")
	flagSet.String("oauth2.provider_config.user_info_url", "", "oauth2 provider user info url")
	flagSet.String("oauth2.provider_config.jwks_url", "", "oauth2 provider jwks url")
	flagSet.StringSlice("oauth2.provider_config.algorithms", nil, "oauth2 provider algorithms")
	common.MustNoError(viper.BindPFlags(flagSet))
	configuration.Register(c)
}

func (c *config) Read() {
	common.MustDecodeFromMapstructure(viper.AllSettings()["oauth2"], c)
}
