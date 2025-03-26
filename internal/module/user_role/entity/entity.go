package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	ID                  uuid.UUID  `db:"id"`
	UserID              uuid.UUID  `db:"user_id"`
	RoleID              uuid.UUID  `db:"role_id"`
	CreatedAt           time.Time  `db:"created_at"`
	DeletedAt           *time.Time `db:"deleted_at"`
	RoleName            string     `db:"role_name"`
	RoleAppPermissionID uuid.UUID  `db:"role_app_permission_id"`
	AppPermissionID     uuid.UUID  `db:"app_permission_id"`
	ApplicationID       uuid.UUID  `db:"application_id"`
	PermissionID        uuid.UUID  `db:"permission_id"`
	Resource            string     `db:"resource"`
	Action              string     `db:"action"`
}
