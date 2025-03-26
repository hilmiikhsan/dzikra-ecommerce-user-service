package repository

import (
	"context"
	"encoding/json"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ApplicationRepository = &applicationRepository{}

type applicationRepository struct {
	db *sqlx.DB
}

func NewApplicationRepository(db *sqlx.DB) *applicationRepository {
	return &applicationRepository{
		db: db,
	}
}

func (r *applicationRepository) FindAllApplication(ctx context.Context) ([]entity.Application, error) {
	var res []entity.Application

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(queryFindAllApplication))
	if err != nil {
		log.Error().Err(err).Any("res", &res).Msg("repository::FindAllApplication - Failed to get list application")
		return nil, err
	}

	return res, nil
}

func (r *applicationRepository) FindPermissionAppsByIDs(ctx context.Context, appIDs []string) ([]dto.PermissionApp, error) {
	query, args, err := sqlx.In(queryGetListPermissionByAppFiltered, appIDs)
	if err != nil {
		log.Error().Err(err).Msg("repository::FindPermissionAppsByIDs - error building query")
		return nil, err
	}
	query = r.db.Rebind(query)

	var rows []entity.ApplicationPermission
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		log.Error().Err(err).Msg("repository::FindPermissionAppsByIDs - error executing query")
		return nil, err
	}

	var result []dto.PermissionApp
	for _, row := range rows {
		var perms []dto.ListPermissionApp
		if err := json.Unmarshal([]byte(row.Permissions), &perms); err != nil {
			log.Error().Err(err).Msg("repository::FindPermissionAppsByIDs - error unmarshalling permissions JSON")
			return nil, err
		}
		result = append(result, dto.PermissionApp{
			ID:          row.ID,
			Name:        row.Name,
			Permissions: perms,
		})
	}

	return result, nil
}
