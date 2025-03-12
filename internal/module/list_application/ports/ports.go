package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/list_application/entity"
)

type ApplicationRepository interface {
	FindAllApplication(ctx context.Context) ([]entity.Application, error)
}
