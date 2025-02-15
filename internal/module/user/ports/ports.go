package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/entity"
)

type UserRepository interface {
	InsertNewUser(ctx context.Context, tx *sql.Tx, data *entity.User) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type UserService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
}
