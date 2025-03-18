package repository

const (
	queryFindApplicationPermissionByIDs = `
		SELECT COUNT(*) FROM application_permissions WHERE id IN (?)
	`

	queryFindApplicationPermissionByActionAndResource = `
		SELECT 
			ap.id
		FROM application_permissions ap
		JOIN permissions p ON ap.permission_id = p.id
		WHERE p.action = ? AND p.resource = ?
	`

	queryGetPermissionIDByAppPermID = `
		SELECT permission_id FROM application_permissions WHERE id = ?
	`
)
