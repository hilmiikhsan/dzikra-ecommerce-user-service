package ports

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user_profile/entity"
)

type UserProfileRepository interface {
	InsertNewUserProfile(ctx context.Context, tx *sql.Tx, data *entity.UserProfile) error
	FindByUserID(ctx context.Context, userID string) (*entity.UserProfile, error)
	SoftDeleteByUserID(ctx context.Context, tx *sql.Tx, userID string) error
}
