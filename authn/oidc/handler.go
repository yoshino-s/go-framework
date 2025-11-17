package authn_oidc

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *OidcAuthentication) LoginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		redirect := c.QueryParam("redirect")
		if redirect != "" {
			c.SetCookie(&http.Cookie{
				Name:     "redirect",
				Value:    redirect,
				HttpOnly: true,
				Expires:  time.Now().Add(5 * time.Minute),
				Path:     "/",
			})
		}
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

func (h *OidcAuthentication) CallbackHandler(enableRedirect bool, handler func(c echo.Context, userId, email string, roles []string) error) echo.HandlerFunc {
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

		if redirectCookie, err := c.Cookie("redirect"); err == nil && enableRedirect {
			// Clear redirect cookie
			c.SetCookie(&http.Cookie{
				Name:     "redirect",
				Value:    "",
				HttpOnly: true,
				Expires:  time.Unix(0, 0),
			})
			return c.Redirect(http.StatusFound, redirectCookie.Value)
		}

		return handler(c, userID, email, roles)
	}
}
