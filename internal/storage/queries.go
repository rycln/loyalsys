package storage

const (
	sqlAddUser = "INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id"

	sqlGetUserByLogin = "SELECT id, login, password_hash FROM users WHERE login = $1"
)
