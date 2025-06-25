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

	GetAllWithdrawDetailsByUser = `
        SELECT o.number, w.sum, w.created_at
        FROM withdraw w
        INNER JOIN "order" o ON w.order_id = o.id
        INNER JOIN user_data u ON o.user_id = u.id
        WHERE u.id = $1
        ORDER BY w.created_at DESC`

	GetTotalAccrualByUser = `
		SELECT COALESCE(sum(o.accrual), 0)
		FROM "order" o 
		WHERE o.status ='PROCESSED' AND o.user_id = $1 
`

	GetTotalWithdrawByUser = `
	    SELECT COALESCE(sum(w.sum), 0)
        FROM withdraw w
        INNER JOIN "order" o ON w.order_id = o.id
        INNER JOIN user_data u ON o.user_id = u.id
        WHERE u.id = $1
`

	GetUnprocessedOrderNumbers = `
		SELECT number
		FROM "order"
		WHERE status IN ('NEW', 'PROCESSING')
`

	UpdateAccrualData = `
	UPDATE "order"
	SET status = $1, accrual = $2
	WHERE number = $3
`
)
