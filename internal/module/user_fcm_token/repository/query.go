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

	queryFindFcmUserTokenByRole = `
		SELECT
    		ufc.fcm_token
		FROM user_fcm_tokens ufc
		JOIN users u ON u.id = ufc.user_id
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE r.name = ?
		AND ufc.fcm_token   IS NOT NULL
		AND ufc.fcm_token  <> ''
		AND ufc.device_id   IS NOT NULL
		AND ufc.device_id  <> ''
		AND ufc.device_type IS NOT NULL  
		AND ufc.deleted_at  IS NULL
		AND u.deleted_at    IS NULL
		AND r.deleted_at    IS NULL
	`
)
