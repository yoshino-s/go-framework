package authentication

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

type Validator interface {
	Validate(ctx echo.Context, token string) (bool, error)
}

type ValidatorFunc func(ctx echo.Context, token string) (bool, error)

func (f ValidatorFunc) Validate(ctx echo.Context, token string) (bool, error) {
	return f(ctx, token)
}

func ConstantValidator(tokens []string) Validator {
	return ValidatorFunc(func(ctx echo.Context, token string) (bool, error) {
		return slices.Contains(tokens, token), nil
	})
}
