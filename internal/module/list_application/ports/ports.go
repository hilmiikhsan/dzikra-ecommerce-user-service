package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/entity"
)

type ApplicationRepository interface {
	FindAllApplication(ctx context.Context) ([]entity.Application, error)
	FindPermissionAppsByIDs(ctx context.Context, appIDs []string) ([]dto.PermissionApp, error)
}

type ApplicationService interface {
	GetListApplication(ctx context.Context) ([]dto.GetListApplicationResponse, error)
	GetListPermissionByApp(ctx context.Context, appIDsParam string) (*dto.GetListPermissionByAppResponse, error)
}
