package hashutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNeUUID(t *testing.T) {
	uuid := NewUUID()
	require.NotZero(t, uuid, "uuid invalid instance")
}

func TestSha1(t *testing.T) {
	hash := Sha1("any string value")
	require.NotZero(t, hash, "sha1 invalid instance")
}

func TestSha1f(t *testing.T) {
	hash := Sha1f("temple string: %s to be hashed %v times", "any value", 2)
	require.NotZero(t, hash, "sha1f invalid instance")
}
