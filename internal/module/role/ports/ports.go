package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/entity"
)

type RoleRepository interface {
	FindRoleByName(ctx context.Context, name string) (*entity.Role, error)
}
