package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_app_permission/entity"
)

type RoleAppPermissionRepository interface {
	InsertNewRoleAppPermissions(ctx context.Context, tx *sql.Tx, data []entity.RoleAppPermission) error
	SoftDeleteRoleAppPermissions(ctx context.Context, tx *sql.Tx, roleID string) error
}
