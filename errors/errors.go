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

func New(message interface{}, code ...int) error {
	e := &AppError{}
	if len(code) > 0 {
		e.code = code[0]
	} else {
		e.code = http.StatusInternalServerError
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

func Wrap(err error, msg string, code ...int) error {
	if err == nil {
		return nil
	}
	e := &AppError{
		error:   err,
		message: msg,
	}
	if len(code) > 0 {
		e.code = code[0]
	} else {
		e.code = http.StatusInternalServerError
	}

	return errors.Wrap(e, 1)
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
