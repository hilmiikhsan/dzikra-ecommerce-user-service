package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/entity"
)

type UserRoleRepository interface {
	InsertNewUserRole(ctx context.Context, tx *sql.Tx, data *entity.UserRole) error
	FindByUserID(ctx context.Context, userID string) ([]string, error)
	FindPermissionsByUserID(ctx context.Context, userID string) ([]string, error)
	SoftDeleteUserRolePermissions(ctx context.Context, tx *sql.Tx, roleID string) error
}
