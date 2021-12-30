package error

import (
	"errors"
	"net/http"
)

//Error service error
type Error struct {
	Code int
	Err  error
}

//Error print error
func (e *Error) Error() string {
	return e.Err.Error()
}

//New create new error
func New(code int, err error) *Error {
	if err == nil {
		return &Error{
			Code: code,
			Err:  errors.New(http.StatusText(code)),
		}
	}

	if e, ok := err.(*Error); ok {
		e.Code = code
		return e
	}

	return &Error{
		Code: code,
		Err:  err,
	}
}

//NotFound 404
func NotFound(err error) *Error {
	return New(http.StatusNotFound, err)
}

//BadRequest 400
func BadRequest(err error) *Error {
	return New(http.StatusBadRequest, err)
}

//InternalServerError 500
func InternalServerError(err error) *Error {
	return New(http.StatusInternalServerError, err)
}
