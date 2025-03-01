package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_fcm_token/entity"
)

type UserFCMTokenRepository interface {
	InsertNewUserFCMToken(ctx context.Context, tx *sql.Tx, userFCMToken entity.UserFCMToken) error
	DeleteUserFCMToken(ctx context.Context, userID string) error
}
