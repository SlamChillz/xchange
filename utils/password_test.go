package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordHash(t *testing.T) {
	password := RandomString(8)
	hash1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)
	require.NotEqual(t, password, hash1)
	
	err = CheckPassword(hash1, password)
	require.NoError(t, err)

	wrongPassword := RandomString(8)
	err = CheckPassword(hash1, wrongPassword)
	require.Error(t, err)

	hash2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)
	require.NotEqual(t, hash1, hash2)
}
