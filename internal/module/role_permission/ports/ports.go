package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_permission/entity"
)

type RolePermissionRepository interface {
	GetUserRolePermission(ctx context.Context, roleID []string) ([]entity.UserRolePermission, error)
}
