package repository

const (
	queryFindAllApplication = `
		SELECT
			id,
			name
		FROM applications
		ORDER BY name ASC
	`

	queryGetListPermissionByAppFiltered = `
		SELECT 
			a.id,
			a.name,
			COALESCE(
				json_agg(
					json_build_object(
						'appperm_id', ap.id,
						'action', p.action,
						'resource', p.resource
					)
				) FILTER (WHERE ap.id IS NOT NULL), '[]'
			) AS permissions
		FROM applications a
		LEFT JOIN application_permissions ap ON a.id = ap.application_id
		LEFT JOIN permissions p ON ap.permission_id = p.id
		WHERE a.id IN (?)
		GROUP BY a.id, a.name
		ORDER BY a.name
	`
)
