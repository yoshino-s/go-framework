package authn_jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/utils"
)

type JwtAuthenticationConfig struct {
	Secret []byte        `json:"secret" yaml:"secret" mapstructure:"secret"`
	Ttl    time.Duration `json:"ttl" yaml:"ttl" mapstructure:"ttl"`
}

var _ configuration.Configuration = (*JwtAuthenticationConfig)(nil)

func (c *JwtAuthenticationConfig) Register(flagSet *pflag.FlagSet) {
	flagSet.BytesBase64("authn.jwt.secret", []byte{}, "JWT secret key")
	utils.MustNoError(viper.BindPFlags(flagSet))
	configuration.Register(c)
}

func (c *JwtAuthenticationConfig) Read() {
	authn, ok := viper.AllSettings()["authn"]
	if !ok {
		panic("authn.jwt configuration not found")
	}

	utils.MustDecodeFromMapstructure(authn.(map[string]any)["jwt"], c)

	if len(c.Secret) == 0 {
		panic("authn.jwt.secret must be set")
	}
}

type JwtMiddlewareConfig struct {
	AuthGetter []AuthGetter
	Skipper    middleware.Skipper
	Public     middleware.Skipper
}

type AuthGetter func(c echo.Context) (bool, string)

// Claims represents application JWT claims.
type Claims struct {
	UserID               string   `json:"sub"`
	Email                string   `json:"email"`
	Roles                []string `json:"roles"`
	jwt.RegisteredClaims `json:",inline"`
}
