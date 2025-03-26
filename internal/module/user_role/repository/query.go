package repository

const (
	queryInsertNewUserRole = `
		INSERT INTO user_roles
		(
			id,
			user_id,
			role_id
		) VALUES (?, ?, ?) ON CONFLICT ON CONSTRAINT unique_user_roles DO NOTHING
	`

	queryFindByUserID = `
		SELECT
			role_id
		FROM user_roles
		WHERE user_id = ?
	`

	queryFindPermissionByUserID = `
		SELECT 
			p.action || '|' || p.resource AS permission
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE ur.user_id = ?
	`

	querySoftDeleteUserRolePermissions = `
		UPDATE user_roles
		SET 
			deleted_at = CURRENT_TIMESTAMP
		WHERE role_id = ? AND deleted_at IS NULL
	`

	querySoftDeleteUserRoles = `
		UPDATE user_roles
		SET 
			deleted_at = CURRENT_TIMESTAMP
		WHERE user_id = ? AND deleted_at IS NULL
	`

	queryFindUserRoleDetailsByUserID = `
		SELECT 
			id, user_id, 
			role_id, 
			created_at, 
			deleted_at
        FROM user_roles
        WHERE user_id = ? AND deleted_at IS NULL
	`

	querySoftDeleteUserRolesByIDs = `
		UPDATE user_roles
        SET deleted_at = CURRENT_TIMESTAMP
        WHERE user_id = ? AND role_id IN (?) AND deleted_at IS NULL
	`

	queryFindAllUserRolesByUserID = `
		SELECT 
			id,
			user_id,
			role_id,
			created_at,
			deleted_at 
		FROM user_roles WHERE user_id = ?
	`

	queryFindUserRoleByUserIDAndRoleName = `
		SELECT 
			ur.id,
			ur.user_id,
			ur.role_id,
			ur.created_at,
			ur.deleted_at
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ? AND UPPER(r.name) = ?
`

	queryRestoreUserRole = `
		UPDATE user_roles 
		SET 
			deleted_at = NULL 
		WHERE id = ?
	`

	querySoftDeleteUserRoleByID = `
		UPDATE user_roles 
		SET 
			deleted_at = CURRENT_TIMESTAMP 
		WHERE id = ? AND deleted_at IS NULL
	`
)
