package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrWrongPassword = errors.New("wrong password")

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func CompareHashAndPassword(hashed, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return ErrWrongPassword
	}
	return nil
}
