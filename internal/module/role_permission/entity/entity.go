package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type RolePermission struct {
	ID           uuid.UUID `db:"id"`
	RoleID       uuid.UUID `db:"role_id"`
	PermissionID uuid.UUID `db:"permission_id"`
}

type UserRolePermission struct {
	RoleName        string         `db:"role_name"`
	ApplicationID   sql.NullString `db:"application_id"`
	ApplicationName sql.NullString `db:"application_name"`
	Permission      string         `db:"permission"`
}
