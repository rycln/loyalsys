package models

type Withdrawal struct {
	ID          int64   `json:"-"`
	Order       string  `json:"order"`
	UserID      UserID  `json:"-"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}
