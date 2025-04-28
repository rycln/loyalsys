package models

import "errors"

var ErrInvalidUser = errors.New("invalid user")

type UserID int64

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	if u.Login == "" || u.Password == "" {
		return ErrInvalidUser
	}
	return nil
}

type UserDB struct {
	ID           UserID
	Login        string
	PasswordHash string
}
