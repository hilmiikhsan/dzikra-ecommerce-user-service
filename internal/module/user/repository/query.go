package repository

const (
	queryInsertNewUser = `
		INSERT INTO users
		(
			id,
			email,
			password,
			full_name
		) VALUES (?, ?, ?, ?) RETURNING id, full_name, email
	`

	queryFindUserByEmail = `
		SELECT
			id,
			password,
			email,
			full_name,
			email_verified_at,
			otp_number_verified_at
		FROM users
		WHERE email = ?
	`

	queryUpdateVerificationUserByEmail = `
		UPDATE users
		SET
			email_verified_at = NOW(),
			otp_number_verified_at = NOW()
		WHERE email = ?
		RETURNING email_verified_at
	`

	queryUpdateUserLastLoginAt = `
		UPDATE users
		SET
			last_login_at = NOW()
		WHERE id = ?
	`

	queryFindUserByID = `
		SELECT
			id,
			password,
			email,
			full_name,
			email_verified_at,
			otp_number_verified_at
		FROM users
		WHERE id = ?
	`

	queryUpdatePasswordByEmail = `
		UPDATE users
		SET
			password = ?
		WHERE email = ?
	`

	queryFindAllUser = `
		WITH user_list AS (
			SELECT 
				u.id,
				u.email,
				u.full_name,
				up.phone_number,
				COALESCE(
					json_agg(
						json_build_object(
							'id', ur.id,
							'user_id', ur.user_id,
							'role_id', ur.role_id
						)
					) FILTER (WHERE ur.id IS NOT NULL), '[]'
				) AS user_role,
				CASE WHEN u.email_verified_at IS NOT NULL THEN true ELSE false END AS email_confirmed
			FROM users u
			LEFT JOIN user_profiles up ON u.id = up.user_id
			LEFT JOIN user_roles ur ON u.id = ur.user_id
			WHERE (u.email ILIKE '%' || ? || '%' OR u.full_name ILIKE '%' || ? || '%')
			AND u.deleted_at IS NULL
			GROUP BY u.id, u.email, u.full_name, up.phone_number, u.email_verified_at
			ORDER BY u.full_name
			LIMIT ? OFFSET ?
		)
			
		SELECT * FROM user_list
	`

	queryCountUser = `
		SELECT COUNT(*) FROM users u
		WHERE (u.email ILIKE '%' || ? || '%' OR u.full_name ILIKE '%' || ? || '%')
		AND u.deleted_at IS NULL
	`

	queryFindUserDetailByID = `
		WITH user_detail AS (
			SELECT 
				u.id,
				u.email,
				u.full_name,
				up.phone_number,
				COALESCE(
					json_agg(
						json_build_object(
							'id', ur.id,
							'user_id', ur.user_id,
							'role_id', ur.role_id
						)
					) FILTER (WHERE ur.id IS NOT NULL), '[]'
				) AS user_role,
				CASE WHEN u.email_verified_at IS NOT NULL THEN true ELSE false END AS email_confirmed
			FROM users u
			LEFT JOIN user_profiles up ON u.id = up.user_id
			LEFT JOIN user_roles ur ON u.id = ur.user_id
			WHERE u.id = ? AND u.deleted_at IS NULL
			GROUP BY u.id, u.email, u.full_name, up.phone_number, u.email_verified_at
		)
			
		SELECT * FROM user_detail;
	`

	queryUpdateNewUser = `
		UPDATE users
		SET
			full_name = ?,
			email = ?,
			password = ?,
			email_verified_at = NULL
		WHERE id = ?
		RETURNING id, full_name, email, email_verified_at
	`

	queryUpdateNewUserWithoutEmail = `
		UPDATE users
		SET
			full_name = ?,
			password = ?
		WHERE id = ?
		RETURNING id, full_name, email, email_verified_at
	`
)
