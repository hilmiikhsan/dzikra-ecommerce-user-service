package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user_fcm_token/entity"
)

type UserFCMTokenRepository interface {
	FindUserFCMTokenDetail(ctx context.Context, deviceID, deviceType, userID string) (*entity.UserFCMToken, error)
	InsertNewUserFCMToken(ctx context.Context, tx *sql.Tx, userFCMToken *entity.UserFCMToken) error
	DeleteUserFCMToken(ctx context.Context, tx *sql.Tx, userID string) error
	UpdateUserFCMToken(ctx context.Context, tx *sql.Tx, userFCMToken *entity.UserFCMToken) error
	FindFcmUserTokenByRole(ctx context.Context, role string) ([]string, error)
	FindUserFCMTokenByUserID(ctx context.Context, userID string) (string, error)
}
