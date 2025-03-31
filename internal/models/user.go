package models

type UserID int64

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
