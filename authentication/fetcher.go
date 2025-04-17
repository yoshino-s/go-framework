package authentication

import "github.com/labstack/echo/v4"

type Fetcher interface {
	Fetch(c echo.Context) (string, error)
}

type FetcherFunc func(c echo.Context) (string, error)

func (f FetcherFunc) Fetch(c echo.Context) (string, error) {
	return f(c)
}

func HeaderFetcher(header string) Fetcher {
	return FetcherFunc(func(c echo.Context) (string, error) {
		return c.Request().Header.Get(header), nil
	})
}

func BasicAuthFetcher() Fetcher {
	return FetcherFunc(func(c echo.Context) (string, error) {
		user, password, ok := c.Request().BasicAuth()
		if !ok {
			return "", echo.ErrUnauthorized
		}
		return user + ":" + password, nil
	})
}

func BearerTokenFetcher() Fetcher {
	return FetcherFunc(func(c echo.Context) (string, error) {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return "", echo.ErrUnauthorized
		}
		return authHeader[7:], nil
	})
}

func QueryFetcher(param string) Fetcher {
	return FetcherFunc(func(c echo.Context) (string, error) {
		return c.QueryParam(param), nil
	})
}

func CookieFetcher(name string) Fetcher {
	return FetcherFunc(func(c echo.Context) (string, error) {
		cookie, err := c.Cookie(name)
		if err != nil {
			return "", err
		}
		return cookie.Value, nil
	})
}
