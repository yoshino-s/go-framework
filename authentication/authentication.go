package authentication

import "github.com/labstack/echo/v4"

type Authentication struct {
	validate func(string) (bool, error)
}

func New(validate func(string) (bool, error)) *Authentication {
	return &Authentication{validate: validate}
}

func (auth *Authentication) Middleware(fetcher func(c echo.Context) (string, error)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := fetcher(c)
			if err != nil {
				return err
			}

			if ok, err := auth.validate(token); err != nil {
				return err
			} else if !ok {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}
