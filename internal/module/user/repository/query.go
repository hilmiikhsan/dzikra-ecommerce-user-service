package repository

const (
	queryInsertNewUser = `
		INSERT INTO users
		(
			id,
			username,
			email,
			password,
			full_name
		) VALUES (?, ?, ?, ?, ?) RETURNING id, username, full_name, email
	`

	queryFindUserByEmail = `
		SELECT
			id,
			username,
			email,
			full_name,
			email_verified_at,
			otp_number_verified_at
		FROM users
		WHERE email = ?
	`
)
