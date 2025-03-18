package repository

const (
	querySoftDeleteRoleAppPermissions = `
		UPDATE role_app_permissions
		SET 
			deleted_at = CURRENT_TIMESTAMP
		WHERE role_id = ? AND deleted_at IS NULL
	`
)
