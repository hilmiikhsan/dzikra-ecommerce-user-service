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

	queryFindUserFCMTokenDetail = `
		SELECT
			device_id,
			device_type
		FROM user_fcm_tokens
		WHERE device_id = ? AND device_type = ? AND user_id = ?
	`

	queryUpdateUserFCMToken = `
		UPDATE user_fcm_tokens
		SET
			fcm_token = ?
		WHERE device_id = ? AND device_type = ? AND user_id = ?
	`
)
