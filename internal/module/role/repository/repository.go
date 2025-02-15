package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/ports"
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
