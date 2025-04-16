package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const testPassword = "secret"

var tooBigPassword = string(make([]byte, 100))

func TestHashPassword(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		hash, err := HashPassword(testPassword)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("too big password", func(t *testing.T) {
		_, err := HashPassword(tooBigPassword)
		assert.Error(t, err)
	})
}

func TestCompareHashAndPassword(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		preHash, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
		require.NoError(t, err)
		err = CompareHashAndPassword(string(preHash), testPassword)
		assert.NoError(t, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		preHash, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
		require.NoError(t, err)
		err = CompareHashAndPassword(string(preHash), "wrong_password")
		assert.ErrorIs(t, err, ErrWrongPassword)
	})
}
