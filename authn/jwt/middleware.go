package authn_jwt

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ContextClaimsKey struct{}

func (j *JwtAuthentication) Middleware(config *JwtMiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}

			var tokenStr string
			found := false
			for _, getter := range config.AuthGetter {
				ok, ts := getter(c)
				if ok {
					tokenStr = ts
					found = true
					break
				}
			}

			isPublic := config.Public != nil && config.Public(c)

			if !isPublic && !found {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			var claims *Claims
			var err error
			claims, err = j.ParseJWT(tokenStr)

			if err == nil {
				c.SetRequest(
					c.Request().WithContext(context.WithValue(c.Request().Context(), ContextClaimsKey{}, claims)),
				)
			} else {
				if !isPublic {
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
				}
			}

			return next(c)
		}
	}
}

// GetClaims helper.
func GetClaims(c context.Context) *Claims {
	if claims, ok := c.Value(ContextClaimsKey{}).(*Claims); ok {
		return claims
	}
	return nil
}
