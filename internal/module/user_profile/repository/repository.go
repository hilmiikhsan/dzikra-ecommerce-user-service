package repository

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.UserProfileRepository = &userProfileRepository{}

type userProfileRepository struct {
	db *sqlx.DB
}

func NewUserProfileRepository(db *sqlx.DB) *userProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

func (r *userProfileRepository) InsertNewUserProfile(ctx context.Context, tx *sql.Tx, data *entity.UserProfile) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewUserProfile),
		data.ID,
		data.UserID,
		data.PhoneNumber,
	)
	if err != nil {
		uniqueConstraints := map[string]string{
			"user_profiles_phone_number_key": constants.ErrPhoneNumberAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::InsertNewUserProfile - Failed to insert new user profile")
			return handleErr
		}

		if _, ok := val.(*entity.UserProfile); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewUserProfile - Failed to insert new user profile")
			return nil
		}

		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}

func (r *userProfileRepository) FindByUserID(ctx context.Context, userID string) (*entity.UserProfile, error) {
	var res = new(entity.UserProfile)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindByUserID), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Str("userID", userID).Msg("repository::FindByUserID - User profile not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrUserProfileNotFound))
		}

		log.Error().Err(err).Str("userID", userID).Msg("repository::FindByUserID - Failed to find user profile by user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}
