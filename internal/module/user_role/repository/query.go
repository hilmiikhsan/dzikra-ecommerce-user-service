package repository

const (
	queryInsertNewUserRole = `
		INSERT INTO user_roles
		(
			id,
			user_id,
			role_id
		) VALUES (?, ?, ?)
	`

	queryFindByUserID = `
		SELECT
			role_id
		FROM user_roles
		WHERE user_id = ?
	`

	queryFindUserRolesWithPermissions = `
		SELECT 
			ur.user_id, 
			r.id AS role_id, 
			r.name AS role_name,
		    rap.id AS role_app_permission_id, 
			rap.app_permission_id,
		    ap.application_id,
		    p.id AS permission_id, 
			p.resource, 
			p.action
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		JOIN role_app_permissions rap ON r.id = rap.role_id
		JOIN application_permissions ap ON rap.app_permission_id = ap.id
		JOIN permissions p ON ap.permission_id = p.id
		WHERE ur.user_id = ? AND rap.app_permission_id IN (?)
	`

	queryFindPermissionByUserID = `
		SELECT 
			p.action || '|' || p.resource AS permission
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE ur.user_id = ?
	`
)
