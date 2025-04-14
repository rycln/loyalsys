package models

type Order struct {
	Number string
	UserID UserID
}

type OrderDB struct {
	ID        int64   `json:"-"`
	Number    string  `json:"number"`
	UserID    UserID  `json:"-"`
	Status    string  `json:"status"`
	Accrual   float64 `json:"accrual,omitempty"`
	CreatedAt string  `json:"uploaded_at"`
}

type OrderAccrual struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
