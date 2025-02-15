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
)
