package password

import (
	"golang.org/x/crypto/bcrypt"
)

type BCryptHasher struct{}

func NewBCryptHasher() *BCryptHasher {
	return &BCryptHasher{}
}

func (h *BCryptHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func (h *BCryptHasher) Compare(hashed, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return newErrWrongPassword(ErrWrongPassword)
	}
	return nil
}
