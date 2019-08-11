package server

import "fmt"

type httpError struct {
	code  int
	cause error
	msg   string
}

func wrapError(code int, err error, format string, args ...interface{}) error {
	return &httpError{
		code:  code,
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

func (e *httpError) Error() string {
	if e.msg != "" && e.cause != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.cause.Error())
	} else if e.cause == nil {
		return e.msg
	} else {
		return e.cause.Error()
	}
}

func (e *httpError) Code() int {
	return e.code
}

func (e *httpError) Cause() error {
	return e.cause
}
