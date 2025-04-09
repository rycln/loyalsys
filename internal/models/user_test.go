package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		u := &User{
			Login:    "not empty",
			Password: "not empty",
		}
		err := u.Validate()
		assert.NoError(t, err)
	})

	t.Run("empty login", func(t *testing.T) {
		u := &User{
			Login:    "",
			Password: "not empty",
		}
		err := u.Validate()
		assert.ErrorIs(t, err, ErrInvalidUser)
	})

	t.Run("empty password", func(t *testing.T) {
		u := &User{
			Login:    "not empty",
			Password: "",
		}
		err := u.Validate()
		assert.ErrorIs(t, err, ErrInvalidUser)
	})
}
