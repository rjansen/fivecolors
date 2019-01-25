package errors

import (
	"github.com/pkg/errors"
)

func New(m string) error {
	return errors.New(m)
}

func Errorf(f string, a ...interface{}) error {
	return errors.Errorf(f, a...)
}

func Cause(e error) error {
	return errors.Cause(e)
}

func Wrap(c error, m string) error {
	if c == nil {
		return New(m)
	}
	return errors.Wrap(c, m)
}

func Wrapf(c error, f string, a ...interface{}) error {
	if c == nil {
		return Errorf(f, a...)
	}
	return Wrapf(c, f, a...)
}

func ErrorWrap(c error, e error) error {
	if c == nil && e == nil {
		return nil
	}
	if e == nil {
		return c
	}
	if c == nil {
		return New(e.Error())
	}
	return Wrap(c, e.Error())
}

func Sum(errs ...error) error {
	var resultErr error
	for _, e := range errs {
		if e == nil {
			continue
		}
		resultErr = ErrorWrap(resultErr, e)
	}
	return resultErr
}
