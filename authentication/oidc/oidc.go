package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/errors"
	"golang.org/x/oauth2"
)

type OIDCAuthentication struct {
	application.EmptyApplication
	config   config
	provider *oidc.Provider
}

func New() *OIDCAuthentication {
	return &OIDCAuthentication{}
}

func (h *OIDCAuthentication) Configuration() configuration.Configuration {
	return &h.config
}

func (h *OIDCAuthentication) Setup(ctx context.Context) {
	if h.config.ProviderConfig.IssuerURL == "" {
		provider := common.Must(oidc.NewProvider(context.TODO(), h.config.IssuerURL))
		h.provider = provider
	} else {
		pc := &oidc.ProviderConfig{
			IssuerURL:     h.config.ProviderConfig.IssuerURL,
			AuthURL:       h.config.ProviderConfig.AuthURL,
			TokenURL:      h.config.ProviderConfig.TokenURL,
			DeviceAuthURL: h.config.ProviderConfig.DeviceAuthURL,
			UserInfoURL:   h.config.ProviderConfig.UserInfoURL,
			JWKSURL:       h.config.ProviderConfig.JWKSURL,
			Algorithms:    h.config.ProviderConfig.Algorithms,
		}
		provider := pc.NewProvider(context.TODO())
		h.provider = provider
	}
	h.Logger.Info("oidc provider is initialized")
}

func (c *OIDCAuthentication) getOauth2Config(redirectURL string) *oauth2.Config {
	if c.provider == nil {
		panic("provider is not initialized")
	}
	return &oauth2.Config{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     c.provider.Endpoint(),
		Scopes:       c.config.Scopes,
	}
}

type MiddlewareConfig struct {
	ExternalURL  string
	RedirectPath string
	CallbackPath string
	PostProcess  PostProcessFunc
}

type PostProcessFunc func(c echo.Context, token *oauth2.Token, user *oidc.UserInfo) error
type RegisterFunc func(app *echo.Echo)

func (h *OIDCAuthentication) Register(config MiddlewareConfig) (RegisterFunc, error) {
	if config.RedirectPath == "" {
		config.RedirectPath = "/api/auth/redirect"
	}
	if config.CallbackPath == "" {
		config.CallbackPath = "/api/auth/callback"
	}

	callbackURL, err := url.JoinPath(config.ExternalURL, config.CallbackPath)
	if err != nil {
		return nil, err
	}

	return func(app *echo.Echo) {
		app.GET(config.RedirectPath, func(c echo.Context) error {
			redirect := c.QueryParam("redirect")
			state, err := randString(16)
			if err != nil {
				return errors.New("internal server error", http.StatusInternalServerError)
			}
			c.SetCookie(&http.Cookie{
				Name:  "state",
				Value: state,
			})
			if redirect != "" {
				c.SetCookie(&http.Cookie{
					Name:  "redirect",
					Value: redirect,
				})
			}

			return c.Redirect(http.StatusFound, h.getOauth2Config(callbackURL).AuthCodeURL(state))
		})

		app.GET(config.CallbackPath, func(c echo.Context) error {
			state, err := c.Cookie("state")
			if err != nil {
				return errors.Wrap(err, http.StatusForbidden)
			}
			if state.Value != c.QueryParam("state") {
				return errors.New("invalid state", http.StatusForbidden)
			}

			oauth2Token, err := h.getOauth2Config(callbackURL).Exchange(c.Request().Context(), c.QueryParam("code"))
			if err != nil {
				return errors.Wrap(err, http.StatusInternalServerError)
			}

			userInfo, err := h.provider.UserInfo(c.Request().Context(), oauth2.StaticTokenSource(oauth2Token))

			if err != nil {
				return errors.New("failed to get user info", http.StatusInternalServerError)
			}

			c.SetCookie(&http.Cookie{
				Name:   "state",
				Value:  "",
				MaxAge: -1,
			})

			var redirectUrl string
			if redirect, err := c.Cookie("redirect"); err == nil {
				redirectUrl = redirect.Value
				c.SetCookie(&http.Cookie{
					Name:   "redirect",
					Value:  "",
					MaxAge: -1,
				})
			} else {
				redirectUrl = "/"
			}

			if config.PostProcess != nil {
				config.PostProcess(c, oauth2Token, userInfo)
			}

			return c.Redirect(http.StatusFound, redirectUrl)
		})
	}, nil
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
