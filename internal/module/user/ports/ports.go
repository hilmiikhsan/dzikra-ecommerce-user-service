package ports

import (
	"context"
	"database/sql"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/entity"
)

type UserRepository interface {
	InsertNewUser(ctx context.Context, tx *sql.Tx, data *entity.User) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateVerificationUserByEmail(ctx context.Context, email string) (time.Time, error)
	UpdateUserLastLoginAt(ctx context.Context, tx *sql.Tx, userID string) error
}

type UserService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Verification(ctx context.Context, req *dto.VerificationRequest) (*dto.VerificationResponse, error)
	SendOtpNumberVerification(ctx context.Context, req *dto.SendOtpNumberVerificationRequest) (*dto.SendOtpNumberVerificationResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}
