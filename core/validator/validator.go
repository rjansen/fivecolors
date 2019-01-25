package validator

import (
	"strconv"
	"strings"

	"github.com/rjansen/fivecolors/core/errors"
)

var (
	ErrValidate = errors.New("ErrValidate")
)

type Validator func() error

func ValidateAll(validators ...Validator) Validator {
	return func() error {
		return Validate(validators...)
	}
}

func Validate(validators ...Validator) error {
	var validateErr error
	for _, v := range validators {
		validateErr = errors.ErrorWrap(validateErr, v())
	}
	return validateErr
}

func ValidateIsIn(s string, values ...string) Validator {
	return func() error {
		return IsIn(s, values...)
	}
}

func IsIn(s string, values ...string) error {
	for _, v := range values {
		if s == v {
			return nil
		}
	}
	return errors.Errorf("ErrValueNotIn:%s", values)
}

func ValidateIsBlank(s string) Validator {
	return func() error {
		return IsBlank(s)
	}
}

func IsBlank(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("ErrIsBlank")
	}
	return nil
}

func IsNumber(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	return err
}
