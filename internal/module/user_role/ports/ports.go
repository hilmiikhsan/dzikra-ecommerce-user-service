package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user_role/entity"
)

type UserRoleRepository interface {
	InsertNewUserRole(ctx context.Context, tx *sql.Tx, data *entity.UserRole) error
	FindByUserID(ctx context.Context, userID string) ([]string, error)
	FindPermissionsByUserID(ctx context.Context, userID string) ([]string, error)
	SoftDeleteUserRolePermissions(ctx context.Context, tx *sql.Tx, roleID string) error
	SoftDeleteUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID string) error
	FindUserRoleDetailsByUserID(ctx context.Context, userID string) ([]entity.UserRole, error)
	SoftDeleteUserRolesByIDs(ctx context.Context, tx *sql.Tx, userID string, roleIDs []string) error
	FindAllUserRolesByUserID(ctx context.Context, userID string) ([]entity.UserRole, error)
	FindUserRoleByUserIDAndRoleName(ctx context.Context, userID, roleName string) (entity.UserRole, bool, error)
	RestoreUserRole(ctx context.Context, tx *sql.Tx, roleUserID string) error
	SoftDeleteUserRoleByID(ctx context.Context, tx *sql.Tx, roleUserID string) error
}
