package repository

const (
	queryInsertNewUserRole = `
		INSERT INTO user_roles
		(
			id,
			user_id,
			role_id
		) VALUES (?, ?, ?)
	`
)
