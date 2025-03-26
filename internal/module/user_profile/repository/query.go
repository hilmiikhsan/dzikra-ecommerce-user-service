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
		WHERE user_id = ? AND deleted_at IS NULL
	`

	querySoftDeleteByUserID = `
		UPDATE user_profiles
		SET 
			deleted_at = NOW()
		WHERE user_id = ? AND deleted_at IS NULL
	`
)
