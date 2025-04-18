package models

type Balance struct {
	UserID    UserID  `json:"-"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
