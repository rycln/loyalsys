package models

type UserID int64

type User struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserDB struct {
	ID           UserID
	Login        string
	PasswordHash string
}
