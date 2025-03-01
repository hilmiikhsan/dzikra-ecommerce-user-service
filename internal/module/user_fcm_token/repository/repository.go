package repository

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_fcm_token/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_fcm_token/ports"
	"github.com/jmoiron/sqlx"
)

var _ ports.UserFCMTokenRepository = &userFcmTokenRepository{}

type userFcmTokenRepository struct {
	db *sqlx.DB
}

func NewUserFcmTokenRepository(db *sqlx.DB) *userFcmTokenRepository {
	return &userFcmTokenRepository{
		db: db,
	}
}

func (r *userFcmTokenRepository) InsertNewUserFCMToken(ctx context.Context, tx *sql.Tx, userFCMToken entity.UserFCMToken) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewUserFCMToken),
		userFCMToken.ID,
		userFCMToken.UserID,
		userFCMToken.DeviceID,
		userFCMToken.DeviceType,
		userFCMToken.FcmToken,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", userFCMToken).Msg("repository::InsertNewUserFCMToken - Failed to insert new user fcm token")
		return err
	}

	return nil
}

func (R *userFcmTokenRepository) DeleteUserFCMToken(ctx context.Context, userID string) error {
	_, err := R.db.ExecContext(ctx, R.db.Rebind(queryDeleteUserFCMToken), userID)
	if err != nil {
		log.Error().Err(err).Str("userID", userID).Msg("repository::DeleteUserFCMToken - Failed to delete user fcm token")
		return err
	}

	return nil
}
