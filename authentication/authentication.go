package authentication

import (
	"slices"

	"github.com/labstack/echo/v4"
)

type Authentication struct {
	validator Validator
	fetcher   Fetcher

	ignorePaths []string
}

func New(validate Validator, fetcher Fetcher, options ...Options) *Authentication {
	auth := &Authentication{
		validator:   validate,
		fetcher:     fetcher,
		ignorePaths: []string{},
	}

	for _, option := range options {
		option(auth)
	}

	return auth
}

func (auth *Authentication) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if slices.Contains(auth.ignorePaths, c.Path()) {
				return next(c)
			}

			token, err := auth.fetcher.Fetch(c)
			if err != nil {
				return err
			}

			if ok, err := auth.validator.Validate(c, token); err != nil {
				return err
			} else if !ok {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}
