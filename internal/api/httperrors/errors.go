package httperrors

import (
	"errors"
	"net/http"
)

type HTTPError struct {
	Err        error
	StatusCode int
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

func New(statusCode int, err error) HTTPError {
	return HTTPError{
		StatusCode: statusCode,
		Err:        err,
	}
}

func NewValidationError(message string) HTTPError {
	return HTTPError{
		StatusCode: http.StatusBadRequest,
		Err:        errors.New(message),
	}
}

func NewInternalServerError(err error) HTTPError {
	return HTTPError{
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
