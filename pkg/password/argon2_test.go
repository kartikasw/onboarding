package password

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordHashAndCheck(t *testing.T) {
	pw := "super-secret-123"

	hash, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// Validate standard Argon2id format
	parts := strings.Split(hash, "$")
	require.Len(t, parts, 6, "hash should split into 6 parts")
	require.Equal(t, "argon2id", parts[1], "algorithm should be argon2id")

	// Correct password
	require.NoError(t, CheckPassword(pw, hash))

	// Wrong password
	require.Error(t, CheckPassword("wrong-password", hash))
}
