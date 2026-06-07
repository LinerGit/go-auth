package password_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/LinerGit/go-auth/internal/password"
)

func TestHash(t *testing.T) {

	service := password.New(4)

	hash, err := service.Hash(
		"password123",
	)

	require.NoError(t, err)

	require.NotEmpty(t, hash)

	require.NotEqual(
		t,
		"password123",
		hash,
	)
}

func TestCompare_Success(
	t *testing.T,
) {

	service := password.New(4)

	hash, err := service.Hash(
		"password123",
	)

	require.NoError(t, err)

	err = service.Compare(
		hash,
		"password123",
	)

	require.NoError(t, err)
}

func TestCompare_InvalidPassword(
	t *testing.T,
) {

	service := password.New(4)

	hash, err := service.Hash(
		"password123",
	)

	require.NoError(t, err)

	err = service.Compare(
		hash,
		"wrong_password",
	)

	require.Error(t, err)
}
