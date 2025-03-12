package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/list_application/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
)

type ApplicationRepository interface {
	FindAllApplication(ctx context.Context) ([]entity.Application, error)
	FindPermissionAppsByIDs(ctx context.Context, appIDs []string) ([]dto.PermissionApp, error)
}
