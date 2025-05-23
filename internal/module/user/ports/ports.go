package ports

import (
	"context"
	"database/sql"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user/entity"
)

type UserRepository interface {
	InsertNewUser(ctx context.Context, tx *sql.Tx, data *entity.User) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateVerificationUserByEmail(ctx context.Context, email string) (time.Time, error)
	UpdateUserLastLoginAt(ctx context.Context, tx *sql.Tx, userID string) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	UpdatePasswordByEmail(ctx context.Context, email, password string) error
	FindAllUser(ctx context.Context, limit, offset int, search string) ([]dto.GetDetailUserResponse, int, error)
	FindUserDetailByID(ctx context.Context, id string) (*dto.GetDetailUserResponse, error)
	UpdateNewUser(ctx context.Context, tx *sql.Tx, data *entity.User, existingEmail string) (*entity.User, error)
	SoftDeleteUserByID(ctx context.Context, tx *sql.Tx, id string) error
}

type UserService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Verification(ctx context.Context, req *dto.VerificationRequest) (*dto.VerificationResponse, error)
	SendOtpNumberVerification(ctx context.Context, req *dto.SendOtpNumberVerificationRequest) (*dto.SendOtpNumberVerificationResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthUserResponse, error)
	Logout(ctx context.Context, accessToken string, locals *middleware.Locals) error
	GetCurrentUser(ctx context.Context, locals *middleware.Locals) (*dto.GetCurrentUserResponse, error)
	RefreshToken(ctx context.Context, accessToken string, locals *middleware.Locals) (*dto.AuthUserResponse, error)
	ForgotPassword(ctx context.Context, req *dto.SendOtpNumberVerificationRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
	GetListUser(ctx context.Context, page, limit int, search string) (*dto.GetListUserResponse, error)
	GetDetailUser(ctx context.Context, userID string) (*dto.GetDetailUserResponse, error)
	CreateUser(ctx context.Context, req *dto.CreateOrUpdateUserRequest) (*dto.CreateUserResponse, error)
	UpdateUser(ctx context.Context, userID string, req *dto.CreateOrUpdateUserRequest) (*dto.UpdateUserResponse, error)
	RemoveUser(ctx context.Context, id string) error
}
