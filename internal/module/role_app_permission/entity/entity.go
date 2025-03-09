package entity

import "github.com/google/uuid"

type RoleAppPermission struct {
	ID              uuid.UUID `db:"id"`
	RoleID          uuid.UUID `db:"role_id"`
	AppPermissionID uuid.UUID `db:"app_permission_id"`
}
