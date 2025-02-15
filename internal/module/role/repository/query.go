package repository

const (
	queryFindRoleByName = `
		SELECT id, name FROM roles WHERE name = ?
	`
)
