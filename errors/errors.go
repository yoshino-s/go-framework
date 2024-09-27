package errors

import (
	"errors"
	"net/http"
	"runtime/debug"
)

type AppError struct {
	trace []byte
	err   error
	code  int
}

var (
	ErrInternal = errors.New("internal error")
)

func NewMissingComponentError(component string) *AppError {
	return New(component+" is missing", http.StatusInternalServerError)
}

// New returns new app error that formats as the given text.
func New(message string, code int) *AppError {
	return PopStack(newAppError(errors.New(message), code))
}

func Wrap(cause error, code int) *AppError {
	causeStackTracer := new(StackTracer)
	if errors.As(cause, causeStackTracer) {
		// If our cause has set a stack trace, and that trace is a child of our own function
		// as inferred by prefix matching our current program counter stack, then we only want
		// to decorate the error message rather than add a redundant stack trace.
		if ancestorOfCause(callers(1), (*causeStackTracer).StackTrace()) {
			return newAppError(cause, code)
		}
	}

	return PopStack(newAppError(cause, code))
}

func newAppError(err error, code int) *AppError {
	if err == nil {
		err = ErrInternal
	}

	return &AppError{
		err:   err,
		trace: debug.Stack(),
		code:  code,
	}
}

// Error returns the string representation of the error message.
func (e *AppError) Error() string {
	return e.err.Error()
}

func (e *AppError) Unwrap() error {
	return e.err
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Extend(err error) *AppError {
	return &AppError{
		err:   errors.Join(e.err, err),
		trace: e.trace,
		code:  e.code,
	}
}
