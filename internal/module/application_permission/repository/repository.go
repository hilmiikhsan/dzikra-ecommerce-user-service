package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/application_permission/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/application_permission/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ApplicationPermissionRepository = &applicationPermissionRepository{}

type applicationPermissionRepository struct {
	db *sqlx.DB
}

func NewApplicationPermissionRepository(db *sqlx.DB) *applicationPermissionRepository {
	return &applicationPermissionRepository{
		db: db,
	}
}

func (r *applicationPermissionRepository) FindApplicationPermissionByID(ctx context.Context, ids []string) (bool, error) {
	if len(ids) == 0 {
		return false, nil
	}

	var count int

	query, args, err := sqlx.In(queryFindApplicationPermissionByIDs, ids)
	if err != nil {
		log.Error().Err(err).Any("ids", ids).Msg("repository::CheckApplicationPermissionsExist - Failed to get application permissions")
		return false, err
	}
	query = r.db.Rebind(query)

	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		log.Error().Err(err).Any("ids", ids).Msg("repository::CheckApplicationPermissionsExist - Failed to rebind application permissions")
		return false, err
	}

	return count == len(ids), nil
}

func (r *applicationPermissionRepository) FindByActionAndResource(ctx context.Context, action, resource string) ([]entity.AppPermission, error) {
	var appPermissions []entity.AppPermission

	if err := r.db.SelectContext(ctx, &appPermissions, r.db.Rebind(queryFindApplicationPermissionByActionAndResource), action, resource); err != nil {
		log.Error().Err(err).Msg("repository::FindByActionAndResource - error executing query")
		return nil, err
	}

	return appPermissions, nil
}

func (r *applicationPermissionRepository) GetPermissionIDByID(ctx context.Context, appPermID string) (string, error) {
	var permissionID string

	err := r.db.GetContext(ctx, &permissionID, r.db.Rebind(queryGetPermissionIDByAppPermID), appPermID)
	if err != nil {
		log.Error().Err(err).Msg("repository::GetPermissionIDByID - error fetching permission_id")
		return "", err
	}

	return permissionID, nil
}
