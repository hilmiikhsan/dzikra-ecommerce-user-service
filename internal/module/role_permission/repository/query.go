package repository

const (
	queryGetUserRolePermission = `
		SELECT
			rl.name as role_name,
			a.id AS application_id,
			a.name AS application_name,
			p.action || '|' || p.resource AS permission
		FROM role_permissions rp
		JOIN roles rl ON rp.role_id = rl.id
		JOIN permissions p ON rp.permission_id = p.id
		JOIN application_permissions ap ON p.id = ap.permission_id
		JOIN applications a ON ap.application_id = a.id
		WHERE rp.role_id IN (?)
	`
)
