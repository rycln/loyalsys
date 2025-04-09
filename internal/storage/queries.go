package storage

const (
	sqlAddUser = "INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id"

	sqlGetUserByLogin = "SELECT id, login, password_hash FROM users WHERE login = $1"

	sqlGetOrderByNum = "SELECT id, number, user_id, status, accrual, created_at FROM orders WHERE number = $1"

	sqlAddOrder = "INSERT INTO orders (number, user_id) VALUES ($1, $2)"
)
