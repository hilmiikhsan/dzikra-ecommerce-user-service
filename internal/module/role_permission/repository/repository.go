package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role_permission/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role_permission/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.RolePermissionRepository = &rolePermissionRepository{}

type rolePermissionRepository struct {
	db *sqlx.DB
}

func NewRolePermissionRepository(db *sqlx.DB) *rolePermissionRepository {
	return &rolePermissionRepository{
		db: db,
	}
}

func (r *rolePermissionRepository) GetUserRolePermission(ctx context.Context, roleIDs []string) ([]entity.UserRolePermission, error) {
	var res []entity.UserRolePermission

	query, args, err := sqlx.In(queryGetUserRolePermission, roleIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		log.Error().Err(err).Any("roleIDs", roleIDs).Msg("repository::GetUserRolePermission - Failed to get user role permission")
		return nil, err
	}

	return res, nil
}

func (r *rolePermissionRepository) SoftDeleteRolePermissions(ctx context.Context, tx *sql.Tx, roleID string) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteRolePermissions), roleID)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("repository::SoftDeleteRolePermissions - Failed to soft delete role permissions")
		return err
	}

	return nil
}

func (r *rolePermissionRepository) InsertNewRolePermissions(ctx context.Context, tx *sql.Tx, data []entity.RolePermission) error {
	if len(data) == 0 {
		return nil
	}

	baseQuery := "INSERT INTO role_permissions (id, role_id, permission_id) VALUES "
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*3)

	for _, rp := range data {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, rp.ID, rp.RoleID, rp.PermissionID)
	}

	query := baseQuery + strings.Join(valueStrings, ", ")
	query = r.db.Rebind(query)

	_, err := tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewRolePermissions - Failed to insert new role permissions")
		return err
	}

	return nil
}
