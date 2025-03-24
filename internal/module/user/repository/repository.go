package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/ports"
	userRole "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/dto"
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

func (r *userRepository) UpdatePasswordByEmail(ctx context.Context, email, password string) error {
	_, err := r.db.ExecContext(ctx, r.db.Rebind(queryUpdatePasswordByEmail), password, email)
	if err != nil {
		log.Error().Err(err).Any("email", email).Msg("repository::UpdatePasswordByEmail - Failed to update password by email")
		return err
	}

	return nil
}

func (r *userRepository) FindAllUser(ctx context.Context, limit, offset int, search string) ([]dto.GetDetailUserResponse, int, error) {
	var rows []entity.ListUserRow

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(queryFindAllUser), search, search, limit, offset); err != nil {
		log.Error().Err(err).Msg("repository::GetListUser - error executing query")
		return nil, 0, err
	}

	var total int
	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountUser), search, search); err != nil {
		log.Error().Err(err).Msg("repository::GetListUser - error counting users")
		return nil, 0, err
	}

	users := make([]dto.GetDetailUserResponse, 0, len(rows))
	for _, row := range rows {
		var roles []userRole.UserRoleDetail

		if err := json.Unmarshal([]byte(row.UserRole), &roles); err != nil {
			log.Error().Err(err).Msg("repository::GetListUser - error unmarshalling user_role JSON")
			return nil, 0, err
		}

		convertedRoles := ConvertUserRoleDetails(roles)

		userDTO := dto.GetDetailUserResponse{
			ID:             row.ID,
			Email:          row.Email,
			FullName:       row.FullName,
			PhoneNumber:    row.PhoneNumber,
			UserRole:       convertedRoles,
			EmailConfirmed: dto.IsConfirmEmail{IsConfirm: row.EmailConfirmed},
		}

		users = append(users, userDTO)
	}

	return users, total, nil
}

func (r *userRepository) FindUserDetailByID(ctx context.Context, id string) (*dto.GetDetailUserResponse, error) {
	var row entity.ListUserRow

	err := r.db.GetContext(ctx, &row, r.db.Rebind(queryFindUserDetailByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("id", id).Msg("repository::FindUserDetailByID - User not found")
			return nil, errors.New(constants.ErrUserNotFound)
		}

		log.Error().Err(err).Any("id", id).Msg("repository::FindUserDetailByID - Failed to execute query")
		return nil, err
	}

	var roles []userRole.UserRoleDetail
	if err := json.Unmarshal([]byte(row.UserRole), &roles); err != nil {
		log.Error().Err(err).Msg("repository::FindUserDetailByID - error unmarshalling user_role JSON")
		return nil, err
	}

	convertedRoles := ConvertUserRoleDetails(roles)

	response := &dto.GetDetailUserResponse{
		ID:             row.ID,
		Email:          row.Email,
		FullName:       row.FullName,
		PhoneNumber:    row.PhoneNumber,
		UserRole:       convertedRoles,
		EmailConfirmed: dto.IsConfirmEmail{IsConfirm: row.EmailConfirmed},
	}

	return response, nil
}

func ConvertUserRoleDetail(urd userRole.UserRoleDetail) userRole.UserRole {
	return userRole.UserRole(urd)
}

func ConvertUserRoleDetails(details []userRole.UserRoleDetail) []userRole.UserRole {
	var result []userRole.UserRole

	for _, d := range details {
		result = append(result, ConvertUserRoleDetail(d))
	}

	return result
}
