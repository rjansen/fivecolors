package util

import (
	"os"

	"github.com/rjansen/fivecolors/core/validator"
)

func Getenv(k, d string) string {
	v := os.Getenv(k)
	if validator.IsBlank(d) == nil && validator.IsBlank(v) != nil {
		return d
	}
	return v
}
