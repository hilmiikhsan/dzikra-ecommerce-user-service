package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.RoleRepository = &roleRepository{}

type roleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) FindRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	var res = new(entity.Role)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindRoleByName), name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("payload", name).Msg("repository::FindRoleByName - Failed to find role")
			return nil, errors.New(constants.ErrRoleNotFound)
		}

		return nil, err
	}

	return res, nil
}

func (r *roleRepository) InsertNewRole(ctx context.Context, tx *sql.Tx, data *entity.Role) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewRole),
		data.ID,
		data.Name,
		data.Description,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewRole - Failed to insert new role")
		return err
	}

	return nil
}

func (r *roleRepository) FindRolePermission(ctx context.Context, roleID string) (*dto.CreateRolePermissionResponse, error) {
	var rows []entity.RolePermission

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(queryFindRolePermission), roleID); err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("repository::FindRolePermission - Failed to execute query")
		return nil, err
	}

	if len(rows) == 0 {
		log.Error().Str("roleID", roleID).Msg("repository::FindRolePermission - No data found")
		return nil, sql.ErrNoRows
	}

	response := dto.CreateRolePermissionResponse{
		ID:          rows[0].RoleID.String(),
		Roles:       rows[0].RoleName,
		Description: rows[0].Description,
	}

	permMap := make(map[string]dto.RoleAppPermissions)
	for _, row := range rows {
		key := row.RoleAppPermissionID.String()
		perm := dto.Permissions{
			ID:       row.PermissionID.String(),
			Resource: row.Resource,
			Action:   row.Action,
		}

		if existing, ok := permMap[key]; ok {
			existing.ApplicationPermission.Permissions = append(existing.ApplicationPermission.Permissions, perm)
			permMap[key] = existing
		} else {
			permMap[key] = dto.RoleAppPermissions{
				ApplicationPermissionID: key,
				ApplicationPermission: dto.ApplicationPermission{
					ApplicationID: row.ApplicationID.String(),
					Permissions:   []dto.Permissions{perm},
				},
			}
		}
	}

	var roleAppPerms []dto.RoleAppPermissions
	for _, v := range permMap {
		roleAppPerms = append(roleAppPerms, v)
	}
	response.RoleAppPermissions = roleAppPerms

	return &response, nil
}
