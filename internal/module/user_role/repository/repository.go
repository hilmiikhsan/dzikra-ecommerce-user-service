package repository

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.UserRoleRepository = &userRoleRepository{}

type userRoleRepository struct {
	db *sqlx.DB
}

func NewUserRoleRepository(db *sqlx.DB) *userRoleRepository {
	return &userRoleRepository{
		db: db,
	}
}

func (r *userRoleRepository) InsertNewUserRole(ctx context.Context, tx *sql.Tx, data *entity.UserRole) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewUserRole),
		data.ID,
		data.UserID,
		data.RoleID,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewUserRole - Failed to insert user role")
		return err
	}

	return nil
}

func (r *userRoleRepository) FindByUserID(ctx context.Context, userID string) ([]string, error) {
	var res []entity.UserRole

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(queryFindByUserID), userID)
	if err != nil {
		log.Error().Err(err).Str("userID", userID).Msg("repository::FindByUserID - Failed to find user role by user ID")
		return nil, err
	}

	var roleIDs []string
	for _, role := range res {
		roleIDs = append(roleIDs, role.RoleID.String())
	}

	return roleIDs, nil
}
