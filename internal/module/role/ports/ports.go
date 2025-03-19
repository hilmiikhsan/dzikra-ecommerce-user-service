package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/entity"
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

type RoleService interface {
	GetDetailRole(ctx context.Context, roleID string) (*dto.GetDetailRoleResponse, error)
	GetListRole(ctx context.Context, page, limit int, search string) (*dto.GetListRole, error)
	CreateRolePermission(ctx context.Context, req *dto.RolePermissionRequest) (*dto.RolePermissionResponse, error)
	RemoveRolePermission(ctx context.Context, roleID string) error
	UpdateRolePermission(ctx context.Context, req *dto.SoftDeleteRolePermissionRequest, roleID string) (*dto.RolePermissionResponse, error)
}
