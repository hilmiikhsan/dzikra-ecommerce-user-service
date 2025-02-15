package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/entity"
)

type UserProfileRepository interface {
	InsertNewUserProfile(ctx context.Context, tx *sql.Tx, data *entity.UserProfile) error
}
