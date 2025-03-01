package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.UserRepository = &userRepository{}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) InsertNewUser(ctx context.Context, tx *sql.Tx, data *entity.User) (*entity.User, error) {
	var res = new(entity.User)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewUser),
		data.ID,
		data.Username,
		data.Email,
		data.Password,
		data.FullName,
	).Scan(
		&res.ID,
		&res.Username,
		&res.FullName,
		&res.Email,
	)
	if err != nil {
		uniqueConstraints := map[string]string{
			"users_username_key": constants.ErrUsernameAlreadyRegistered,
			"users_email_key":    constants.ErrEmailAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::InsertNewUser - Failed to insert new user")
			return nil, handleErr
		}

		if user, ok := val.(*entity.User); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewUser - Failed to insert new user")
			return user, nil
		}

		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var res = new(entity.User)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindUserByEmail), email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("email", email).Msg("repository::FindUserByEmail - Email not found")
			return nil, nil
		}

		log.Error().Err(err).Any("email", email).Msg("repository::FindUserByEmail - Failed to find user by email")
		return nil, err
	}

	return res, nil
}

func (r *userRepository) UpdateVerificationUserByEmail(ctx context.Context, email string) (time.Time, error) {
	var emailVerifiedAt time.Time

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryUpdateVerificationUserByEmail), email).Scan(&emailVerifiedAt)
	if err != nil {
		log.Error().Err(err).Any("email", email).Msg("repository::UpdateVerificationUserByEmail - Failed to update verification user by email")
		return time.Time{}, err
	}

	return emailVerifiedAt, nil
}

func (r *userRepository) UpdateUserLastLoginAt(ctx context.Context, tx *sql.Tx, userID string) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryUpdateUserLastLoginAt), userID)
	if err != nil {
		log.Error().Err(err).Any("user_id", userID).Msg("repository::UpdateUserLastLoginAt - Failed to update user last login at")
		return err
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var res = new(entity.User)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindUserByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("id", id).Msg("repository::FindUserByID - ID not found")
			return nil, nil
		}

		log.Error().Err(err).Any("id", id).Msg("repository::FindUserByID - Failed to find user by ID")
		return nil, err
	}

	return res, nil
}
