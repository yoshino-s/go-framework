package errors

import (
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
)

type AppError struct {
	error   error
	message string
	code    int
}

func NewMissingComponentError(component string) error {
	return New(component+" is missing", http.StatusInternalServerError)
}

func New(message interface{}, code int) error {
	e := &AppError{
		code: code,
	}
	switch message := message.(type) {
	case error:
		e.error = message
		e.message = message.Error()
	default:
		e.message = fmt.Sprintf("%v", message)
	}

	return errors.Wrap(e, 1)
}

func Wrap(err error, code int) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(&AppError{
		error:   err,
		message: err.Error(),
		code:    code,
	}, 1)
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Error() string {
	if e.error == nil {
		return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
	} else {
		return fmt.Sprintf("code: %d, message: %s, error: %s", e.code, e.message, e.error.Error())
	}
}

func (e *AppError) Unwrap() error {
	return e.error
}
