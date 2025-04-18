package storage

const sqlAddUser = `
	INSERT INTO users (login, password_hash) 
	VALUES ($1, $2) 
	RETURNING id
`

const sqlGetUserByLogin = `
	SELECT 
		id, 
		login, 
		password_hash 
	FROM users 
	WHERE login = $1
`

const sqlGetOrderByNum = `
	SELECT 
		id, 
		number, 
		user_id, 
		status, 
		accrual, 
		created_at 
	FROM orders 
	WHERE number = $1
`

const sqlAddOrder = `
	INSERT INTO orders (number, user_id) 
	VALUES ($1, $2)
`

const sqlGetInconclusiveOrderNums = `
	SELECT 
		number 
	FROM orders 
	WHERE status NOT IN ('INVALID', 'PROCESSED')
`

const sqlUpdateOrdersBatch = `
	UPDATE orders 
	SET 
		status = $1, 
		accrual = $2 
	WHERE number = $3
`

const sqlGetOrdersByUserID = `
	SELECT 
		number, 
		status, 
		accrual, 
		created_at 
	FROM orders 
	WHERE user_id = $1 
	ORDER BY created_at DESC
`

const sqlGetWithdrawalsByUserID = `
	SELECT 
		id, 
		order, 
		sum, 
		processed_at 
	FROM withdrawals 
	WHERE user_id = $1 
	ORDER BY processed_at DESC
`

const sqlGetBalanceByUserID = `
	SELECT 
		COALESCE(SUM(orders.accrual), 0) AS current, 
		COALESCE(SUM(withdrawals.sum), 0) AS withdrawn 
	FROM orders 
	FULL JOIN withdrawals ON orders.user_id = withdrawals.user_id 
	WHERE orders.user_id = $1 OR withdrawals.user_id = $1
`
