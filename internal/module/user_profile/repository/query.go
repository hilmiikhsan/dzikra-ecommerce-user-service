package repository

const (
	queryInsertNewUserProfile = `
		INSERT INTO user_profiles
		(
			id,
			user_id,
			phone_number
		) VALUES (?, ?, ?)
	`

	queryFindByUserID = `
		SELECT
			user_id,
			phone_number
		FROM user_profiles
		WHERE user_id = ?
	`
)
