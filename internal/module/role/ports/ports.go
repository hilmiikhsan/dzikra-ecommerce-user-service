package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
)

type RoleRepository interface {
	FindRoleByName(ctx context.Context, name string) (*entity.Role, error)
	InsertNewRole(ctx context.Context, tx *sql.Tx, data *entity.Role) error
	FindRolePermission(ctx context.Context, roleID string) (*dto.RolePermissionResponse, error)
	FindListRole(ctx context.Context, limit, offset int, search string) ([]dto.GetListRolePermission, int, error)
	FindRoleByID(ctx context.Context, roleID string) (*dto.GetListRolePermission, error)
	SoftDeleteRole(ctx context.Context, tx *sql.Tx, roleID string) error
	UpdateRole(ctx context.Context, tx *sql.Tx, roleID, newName, description, currentName string) error
}
