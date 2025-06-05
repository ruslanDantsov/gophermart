package query

const (
	InsertOrUpdateUserData = `
		INSERT INTO user_data (id, login, password, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET login = $2, password = $3;
	`
)
