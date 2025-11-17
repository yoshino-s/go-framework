package authn_oidc

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *OidcAuthentication) LoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		state, err := randomState(12)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		// Set state cookie (short lived)
		cookie := &http.Cookie{Name: "oauth_state", Value: state, Expires: time.Now().Add(5 * time.Minute), HttpOnly: true, Path: "/"}
		c.SetCookie(cookie)
		return c.Redirect(http.StatusFound, h.authCodeURL(state))
	}
}

func (h *OidcAuthentication) CallbackHandler(handler func(c echo.Context, userId, email string, roles []string) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		state := c.QueryParam("state")
		code := c.QueryParam("code")
		stCookie, err := c.Cookie("oauth_state")
		if err != nil || stCookie.Value != state {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid state")
		}
		_, claims, err := h.exchangeAndVerify(context.Background(), code)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, err.Error())
		}
		userID, email, roles := h.tokenUser(claims)

		return handler(c, userID, email, roles)
	}
}

func (h *OidcAuthentication) LogoutHandler(redirectURL string) echo.HandlerFunc {
	return func(c echo.Context) error {
		var claim struct {
			EndSessionEndpoint string `json:"end_session_endpoint"`
		}

		if err := h.Provider.Claims(&claim); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if claim.EndSessionEndpoint == "" {
			// No end session endpoint, just redirect to home
			return c.Redirect(http.StatusFound, redirectURL)
		}

		u, err := url.Parse(claim.EndSessionEndpoint)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		q := u.Query()
		q.Set("post_logout_redirect_uri", redirectURL)
		u.RawQuery = q.Encode()

		return c.Redirect(http.StatusFound, u.String())
	}
}
