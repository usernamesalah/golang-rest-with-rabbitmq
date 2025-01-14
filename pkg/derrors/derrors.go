// Package derrors defines internal error values to categorize the different
// types error semantics we support.
// inspired by https://pkg.go.dev/golang.org/x/pkgsite/internal/derrors
package derrors

import (
	"errors"
	"fmt"
	"net/http"
)

// Error represents an error that could be wrapping another error, it includes a code for determining what
// triggered the error.
type Error struct {
	orig error
	code ErrorCode
	msg  string
}

// ErrorCode defines supported error status.
type ErrorCode uint

const (
	Unknown ErrorCode = iota
	NotFound
	InvalidArgument
	Duplicate
	Unauthorized
	Forbidden
)

var codes = []struct {
	code   ErrorCode
	status int
}{
	{Unknown, http.StatusInternalServerError},
	{NotFound, http.StatusNotFound},
	{InvalidArgument, http.StatusBadRequest},
	{Duplicate, http.StatusBadRequest},
	{Unauthorized, http.StatusUnauthorized},
	{Forbidden, http.StatusForbidden},
}

// ToStatus returns a status code corresponding to err.
func ToStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var ierr *Error
	if errors.As(err, &ierr) {
		for _, e := range codes {
			if ierr.code == e.code {
				return e.status
			}
		}
	}
	return http.StatusInternalServerError
}

// IsErrCode Compare error code
func IsErrCode(err error, code ErrorCode) bool {
	if err == nil {
		return false
	}
	var ierr *Error
	if errors.As(err, &ierr) {
		if ierr.code == code {
			return true
		}
	}
	return false
}

// Error returns the message, when wrapping errors the wrapped error is returned.
func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.orig)
	}

	return e.msg
}

// New instantiates a new error.
func New(code ErrorCode, format string, args ...interface{}) error {
	return &Error{
		code: code,
		msg:  fmt.Sprintf(format, args...),
	}
}

// WrapStack returns a wrapped error with error code.
func WrapStack(orig error, code ErrorCode, format string, args ...interface{}) error {
	if orig == nil {
		return nil
	}

	var ierr *Error
	if errors.As(orig, &ierr) {
		return Wrap(&orig, format, args...)
	}

	return &Error{
		orig: orig,
		code: code,
		msg:  fmt.Sprintf(format, args...),
	}
}

func Wrap(errp *error, format string, args ...any) error {
	if *errp != nil {
		*errp = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *errp)
	}
	return *errp
}

// Unwrap returns the wrapped error, if any.
func (e *Error) Unwrap() error {
	return e.orig
}

// Code returns the code representing this error.
func (e *Error) Code() ErrorCode {
	return e.code
}
