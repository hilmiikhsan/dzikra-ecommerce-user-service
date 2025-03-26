package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/application_permission/entity"
)

type ApplicationPermissionRepository interface {
	FindApplicationPermissionByID(ctx context.Context, ids []string) (bool, error)
	FindByActionAndResource(ctx context.Context, action, resource string) ([]entity.AppPermission, error)
	GetPermissionIDByID(ctx context.Context, appPermID string) (string, error)
}
