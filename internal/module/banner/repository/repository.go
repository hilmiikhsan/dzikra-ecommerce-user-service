package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.BannerRepository = &bannerRepository{}

type bannerRepository struct {
	db *sqlx.DB
}

func NewBannerRepository(db *sqlx.DB) *bannerRepository {
	return &bannerRepository{
		db: db,
	}
}

func (r *bannerRepository) InsertNewBanner(ctx context.Context, data *entity.Banner) (*entity.Banner, error) {
	var res = new(entity.Banner)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryInsertNewBanner),
		data.ImageURL,
		data.Description,
	).Scan(
		&res.ID,
		&res.ImageURL,
		&res.Description,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewBanner - Failed to insert new banner")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}
