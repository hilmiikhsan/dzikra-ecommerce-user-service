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
          AND r.deleted_at IS NULL
          AND rap.deleted_at IS NULL
	`

	queryFindListRole = `
		WITH role_list AS (
			SELECT 
				r.id,
				r.name AS roles,
				r.description AS "desc",
				COALESCE(
					json_agg(
						json_build_object(
							'application_id', a.id,
							'name', a.name,
							'permission', (
								SELECT json_agg(
									json_build_object(
										'action', p.action,
										'aplicationperm_id', ap.id,
										'resource', p.resource
									)
								)
								FROM application_permissions ap
								JOIN permissions p ON ap.permission_id = p.id
								WHERE ap.id = rap.app_permission_id
							)
						)
					) FILTER (WHERE rap.id IS NOT NULL), '[]'
				) AS role_app_permission
			FROM roles r
			LEFT JOIN role_app_permissions rap ON r.id = rap.role_id
			LEFT JOIN application_permissions ap ON rap.app_permission_id = ap.id
			LEFT JOIN applications a ON ap.application_id = a.id
			WHERE r.name ILIKE '%' || ? || '%' AND r.deleted_at IS NULL
			GROUP BY r.id, r.name, r.description
			ORDER BY r.name
			LIMIT ? OFFSET ?
		)
		
		SELECT * FROM role_list;
	`

	queryCountListRole = `
		SELECT COUNT(*) FROM roles r
		WHERE r.name ILIKE '%' || ? || '%' AND r.deleted_at IS NULL
	`

	queryFindRoleByID = `
		SELECT 
			r.id,
			r.name AS roles,
			r.description AS "desc",
			COALESCE(
				json_agg(
				json_build_object(
					'application_id', a.id,
					'name', a.name,
					'permission', (
					SELECT json_agg(
						json_build_object(
						'action', p.action,
						'aplicationperm_id', ap.id,
						'resource', p.resource
						)
					)
					FROM application_permissions ap
					JOIN permissions p ON ap.permission_id = p.id
					WHERE ap.id = rap.app_permission_id
					)
				)
				) FILTER (WHERE rap.id IS NOT NULL), '[]'
			) AS role_app_permission
		FROM roles r
		LEFT JOIN role_app_permissions rap ON r.id = rap.role_id
		LEFT JOIN application_permissions ap ON rap.app_permission_id = ap.id
		LEFT JOIN applications a ON ap.application_id = a.id
		WHERE r.id = ? AND r.deleted_at IS NULL
		GROUP BY r.id, r.name, r.description
	`

	querySoftDeleteRole = `
		UPDATE roles
		SET 
			deleted_at = CURRENT_TIMESTAMP
		WHERE role_id = ? AND deleted_at IS NULL
	`

	queryUpdateRole = `
		UPDATE roles
		SET 
			name = ?, 
			description = ?, 
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	queryUpdateRoleDescription = `
		UPDATE roles
		SET 
			description = ?, 
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`
)
