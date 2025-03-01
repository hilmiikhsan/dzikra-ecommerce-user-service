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

	queryFindByUserID = `
		SELECT
			role_id
		FROM user_roles
		WHERE user_id = ?
	`
)
