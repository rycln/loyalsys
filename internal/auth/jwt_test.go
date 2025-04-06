package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserID = models.UserID(1)
	testKey    = "secret_key"
)

func TestNewJWTString(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		jwtString, err := NewJWTString(testUserID, testKey)
		assert.NoError(t, err)
		assert.NotEmpty(t, jwtString)
	})
}

func TestParseIDFromJWT(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		claims := jwt.MapClaims{
			"uid": testUserID,
			"exp": tokenExp,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testKey))
		require.NoError(t, err)
		restoredToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(testKey), nil
		})
		require.NoError(t, err)
		uid, err := ParseIDFromJWT(restoredToken)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, uid)
	})

	t.Run("no user id", func(t *testing.T) {
		claims := jwt.MapClaims{
			"exp": tokenExp,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testKey))
		require.NoError(t, err)
		restoredToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(testKey), nil
		})
		require.NoError(t, err)
		_, err = ParseIDFromJWT(restoredToken)
		assert.ErrorIs(t, err, ErrInvalidJWT)
	})
}
