package repository

const (
	querySoftDeleteRoleAppPermissions = `
		UPDATE role_app_permissions
		SET 
			deleted_at = NOW()
		WHERE role_id = ? AND deleted_at IS NULL
	`
)
