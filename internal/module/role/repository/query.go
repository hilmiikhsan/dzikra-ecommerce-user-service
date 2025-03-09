package repository

const (
	queryFindRoleByName = `
		SELECT id, name FROM roles WHERE name = ?
	`

	queryInsertNewRole = `
		INSERT INTO roles 
		(
			id,
			name, 
			description
		) VALUES (?, ?, ?)
	`

	queryFindRolePermission = `
		SELECT
			r.id AS role_id,
			r.name AS role_name,
			r.description,
			rap.id AS role_app_permission_id,
			ap.application_id,
			p.id AS permission_id,
			p.resource,
			p.action
		FROM roles r
		JOIN role_app_permissions rap ON r.id = rap.role_id
		JOIN application_permissions ap ON rap.app_permission_id = ap.id
		JOIN permissions p ON ap.permission_id = p.id
		WHERE r.id = ?
	`
)
