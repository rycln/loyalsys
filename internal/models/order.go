package models

type Order struct {
	Number string
	UserID UserID
}

type OrderDB struct {
	ID        int64
	Number    string
	UserID    UserID
	Status    string
	Accrual   int
	CreatedAt string
}

type OrderAccrual struct {
	Number  string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
