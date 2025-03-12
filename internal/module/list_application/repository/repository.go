package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/list_application/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/list_application/ports"
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
