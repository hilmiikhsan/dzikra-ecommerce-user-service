package repository

const (
	queryInsertNewUserFCMToken = `
		INSERT INTO user_fcm_tokens
		(
			id,
			user_id,
			device_id,
			device_type,
			fcm_token
		) VALUES (?, ?, ?, ?, ?)
	`

	queryDeleteUserFCMToken = `
		DELETE FROM user_fcm_tokens
		WHERE user_id = ?
	`
)
