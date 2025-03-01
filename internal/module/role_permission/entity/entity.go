package entity

import "database/sql"

type UserRolePermission struct {
	RoleName        string         `db:"role_name"`
	ApplicationID   sql.NullString `db:"application_id"`
	ApplicationName sql.NullString `db:"application_name"`
	Permission      string         `db:"permission"`
}
