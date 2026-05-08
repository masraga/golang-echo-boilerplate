package traceerr

import (
	"errors"
	"runtime"
)

type Error struct {
	err  error
	file string
	line int
}

func Wrap(err error) error {
	return wrapWithCaller(err, 2)
}

func WrapReturn(err *error) {
	if err == nil || *err == nil {
		return
	}
	*err = wrapWithCaller(*err, 2)
}

func wrapWithCaller(err error, skip int) error {
	if err == nil {
		return nil
	}
	if _, _, ok := Location(err); ok {
		return err
	}

	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return err
	}

	return &Error{
		err:  err,
		file: file,
		line: line,
	}
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) File() string {
	return e.file
}

func (e *Error) Line() int {
	return e.line
}

func Location(err error) (file string, line int, ok bool) {
	var traced *Error
	if !errors.As(err, &traced) {
		return "", 0, false
	}
	return traced.File(), traced.Line(), true
}
