package repository

const (
	queryFindAllApplication = `
		SELECT
			id,
			name
		FROM applications
		ORDER BY name ASC
	`
)
