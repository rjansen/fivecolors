package hashutil

import (
	"crypto/sha1"
	"fmt"

	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.New().String()
}

func Sha1(v string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(v)))
}

func Sha1f(f string, a ...interface{}) string {
	return Sha1(fmt.Sprintf(f, a...))
}
