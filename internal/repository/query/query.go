package query

const (
	InsertOrUpdateUserData = `
		INSERT INTO user_data (id, login, password, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET login = $2, password = $3;
	`

	FindUserByLogin = `
		SELECT id, login, password, created_at
		FROM user_data
		WHERE login = $1;
`

	InsertOrder = `
		INSERT INTO "order" (id, number, status, accrual, created_at, user_id)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	FindUserByOrderNumber = `
		SELECT user_id 
		FROM "order" 
		WHERE number = $1 FOR UPDATE;
	`
	GetAllOrdersByUser = `
        SELECT id, number, status, accrual, created_at, user_id
        FROM "order" 
        WHERE user_id = $1
        ORDER BY created_at DESC`

	InsertWithdraw = `
		INSERT INTO withdraw (id, sum, created_at, order_id)
		VALUES ($1, $2, $3, $4);
	`
)
