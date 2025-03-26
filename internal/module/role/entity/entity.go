package entity

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
}

type RolePermission struct {
	RoleID              uuid.UUID `db:"role_id"`
	RoleName            string    `db:"role_name"`
	Description         string    `db:"description"`
	RoleAppPermissionID uuid.UUID `db:"role_app_permission_id"`
	ApplicationID       uuid.UUID `db:"application_id"`
	PermissionID        uuid.UUID `db:"permission_id"`
	Resource            string    `db:"resource"`
	Action              string    `db:"action"`
}

type ListRolePermission struct {
	ID                string `db:"id"`
	Roles             string `db:"roles"`
	Description       string `db:"desc"`
	Static            bool   `db:"static"`
	RoleAppPermission string `db:"role_app_permission"`
}
