package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_app_permission/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_app_permission/ports"
	"github.com/jmoiron/sqlx"
)

var _ ports.RoleAppPermissionRepository = &roleAppPermissionRepository{}

type roleAppPermissionRepository struct {
	db *sqlx.DB
}

func NewRoleAppPermissionRepository(db *sqlx.DB) *roleAppPermissionRepository {
	return &roleAppPermissionRepository{
		db: db,
	}
}

func (r *roleAppPermissionRepository) InsertNewRoleAppPermissions(ctx context.Context, tx *sql.Tx, data []entity.RoleAppPermission) error {
	if len(data) == 0 {
		return nil
	}

	baseQuery := "INSERT INTO role_app_permissions (id, role_id, app_permission_id) VALUES "
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*3)

	for _, rp := range data {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, rp.ID, rp.RoleID, rp.AppPermissionID)
	}

	query := baseQuery + strings.Join(valueStrings, ", ")
	query = r.db.Rebind(query)

	_, err := tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewRoleAppPermissions - Failed to insert new role app permissions")
		return err
	}

	return nil
}

func (r *roleAppPermissionRepository) SoftDeleteRoleAppPermissions(ctx context.Context, tx *sql.Tx, roleID string) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteRoleAppPermissions), roleID)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("repository::SoftDeleteRoleAppPermissions - Failed to soft delete role app permissions")
		return err
	}

	return nil
}
