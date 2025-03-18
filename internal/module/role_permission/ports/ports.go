package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_permission/entity"
)

type RolePermissionRepository interface {
	GetUserRolePermission(ctx context.Context, roleID []string) ([]entity.UserRolePermission, error)
	SoftDeleteRolePermissions(ctx context.Context, tx *sql.Tx, roleID string) error
	InsertNewRolePermissions(ctx context.Context, tx *sql.Tx, data []entity.RolePermission) error
}
