package authn_jwt

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	ContextClaimsKey = "claims"
)

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

			if found {
				claims, err := j.ParseJWT(tokenStr)
				if err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
				}
				c.Set(ContextClaimsKey, claims)
			} else {
				if !(config.Public != nil && config.Public(c)) {
					return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
				}
			}
			return next(c)
		}
	}
}

// GetClaims helper.
func GetClaims(c echo.Context) *Claims {
	if v, ok := c.Get(ContextClaimsKey).(*Claims); ok {
		return v
	}
	return nil
}
