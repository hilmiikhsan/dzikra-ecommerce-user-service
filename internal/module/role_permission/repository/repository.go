package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_permission/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_permission/ports"
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
