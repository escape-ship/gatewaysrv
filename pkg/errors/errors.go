package errors

import (
	"fmt"
	"runtime"
)

type Error struct {
	Message string
	Code    string
	Cause   error
	Stack   string
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(message string) *Error {
	return &Error{
		Message: message,
		Stack:   getStack(),
	}
}

func Newf(format string, args ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Stack:   getStack(),
	}
}

func Wrap(err error, message string) *Error {
	return &Error{
		Message: message,
		Cause:   err,
		Stack:   getStack(),
	}
}

func Wrapf(err error, format string, args ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
		Stack:   getStack(),
	}
}

func WithCode(err error, code string) *Error {
	if e, ok := err.(*Error); ok {
		e.Code = code
		return e
	}
	
	return &Error{
		Message: err.Error(),
		Code:    code,
		Cause:   err,
		Stack:   getStack(),
	}
}

func getStack() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}

// Common error codes
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeAuthentication = "AUTHENTICATION_ERROR"
	ErrCodeAuthorization  = "AUTHORIZATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)