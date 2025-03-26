package repository

import (
	"context"
	"database/sql"
	"strings"

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

func (r *userRoleRepository) FindPermissionsByUserID(ctx context.Context, userID string) ([]string, error) {
	var permissions []string

	if err := r.db.SelectContext(ctx, &permissions, r.db.Rebind(queryFindPermissionByUserID), userID); err != nil {
		log.Error().Err(err).Msg("repository::FindPermissionsByUserID - error executing query")
		return nil, err
	}

	return permissions, nil
}

func (r *userRoleRepository) SoftDeleteUserRolePermissions(ctx context.Context, tx *sql.Tx, roleID string) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteUserRolePermissions), roleID)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("repository::SoftDeleteUserRolePermissions - Failed to soft delete user role permissions")
		return err
	}

	return nil
}

func (r *userRoleRepository) SoftDeleteUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID string) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteUserRoles), userID)
	if err != nil {
		log.Error().Err(err).Str("userID", userID).Msg("repository::SoftDeleteUserRoles - Failed to soft delete user roles")
		return err
	}

	return nil
}

func (r *userRoleRepository) FindUserRoleDetailsByUserID(ctx context.Context, userID string) ([]entity.UserRole, error) {
	var roles []entity.UserRole

	if err := r.db.SelectContext(ctx, &roles, r.db.Rebind(queryFindUserRoleDetailsByUserID), userID); err != nil {
		log.Error().Err(err).Str("userID", userID).Msg("repository::FindUserRoleDetailsByUserID - error executing query")
		return nil, err
	}

	return roles, nil
}

func (r *userRoleRepository) SoftDeleteUserRolesByIDs(ctx context.Context, tx *sql.Tx, userID string, roleIDs []string) error {
	if len(roleIDs) == 0 {
		return nil
	}

	query, args, err := sqlx.In(querySoftDeleteUserRolesByIDs, userID, roleIDs)
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteUserRolesByIDs - error building query")
		return err
	}
	query = r.db.Rebind(query)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteUserRolesByIDs - error executing query")
		return err
	}

	return nil
}

func (r *userRoleRepository) FindAllUserRolesByUserID(ctx context.Context, userID string) ([]entity.UserRole, error) {
	var roles []entity.UserRole

	err := r.db.SelectContext(ctx, &roles, r.db.Rebind(queryFindAllUserRolesByUserID), userID)
	if err != nil {
		log.Error().Err(err).Str("userID", userID).Msg("repository::FindAllUserRolesByUserID - Failed to find all user roles by user ID")
		return nil, err
	}

	return roles, nil
}

func (r *userRoleRepository) FindUserRoleByUserIDAndRoleName(ctx context.Context, userID, roleName string) (entity.UserRole, bool, error) {
	var result entity.UserRole

	err := r.db.GetContext(ctx, &result, r.db.Rebind(queryFindUserRoleByUserIDAndRoleName), userID, strings.ToUpper(roleName))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("userID", userID).Any("roleName", roleName).Msg("repository::FindUserRoleByUserIDAndRoleName - User role not found")
			return result, false, nil
		}

		log.Error().Err(err).Any("userID", userID).Any("roleName", roleName).Msg("repository::FindUserRoleByUserIDAndRoleName - Failed to find user role by user ID and role name")
		return result, false, err
	}

	return result, true, nil
}

func (r *userRoleRepository) RestoreUserRole(ctx context.Context, tx *sql.Tx, roleUserID string) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryRestoreUserRole), roleUserID)
	if err != nil {
		log.Error().Err(err).Str("roleUserID", roleUserID).Msg("repository::RestoreUserRole - Failed to restore user role")
		return err
	}

	return nil
}

func (r *userRoleRepository) SoftDeleteUserRoleByID(ctx context.Context, tx *sql.Tx, roleUserID string) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteUserRoleByID), roleUserID)
	if err != nil {
		log.Error().Err(err).Str("roleUserID", roleUserID).Msg("repository::SoftDeleteUserRoleByID - Failed to soft delete user role")
		return err
	}

	return nil
}
