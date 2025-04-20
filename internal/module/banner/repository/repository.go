package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
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

func (r *bannerRepository) FindListBanner(ctx context.Context, limit, offset int, search string) ([]dto.GetListBanner, int, error) {
	var responses []entity.Banner

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListBanner), search, limit, offset); err != nil {
		log.Error().Err(err).Msg("repository::FindListBanner - error executing query")
		return nil, 0, err
	}

	var total int

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountListBanner), search); err != nil {
		log.Error().Err(err).Msg("repository::FindListVoucher - error counting banner")
		return nil, 0, err
	}

	publicURL := config.Envs.MinioStorage.PublicURL

	banners := make([]dto.GetListBanner, 0, len(responses))
	for _, v := range responses {
		banners = append(banners, dto.GetListBanner{
			ID:          v.ID,
			ImageURL:    utils.FormatMediaPathURL(v.ImageURL, publicURL),
			Description: v.Description,
		})
	}

	return banners, total, nil
}

func (r *bannerRepository) UpdateBanner(ctx context.Context, data *entity.Banner) (*entity.Banner, error) {
	var res = new(entity.Banner)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryUpdateBanner),
		data.ImageURL,
		data.Description,
		data.ID,
	).Scan(
		&res.ID,
		&res.ImageURL,
		&res.Description,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateBanner - Failed to update banner")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *bannerRepository) FindBannerByID(ctx context.Context, id int) (*entity.Banner, error) {
	var res = new(entity.Banner)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryFindBannerByID), id).Scan(
		&res.ID,
		&res.ImageURL,
		&res.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("id", id).Msg("repository::FindBannerByID - banner is not found")
			return nil, errors.New(constants.ErrBannerNotFound)
		}

		log.Error().Err(err).Msg("repository::FindBannerByID - Failed to find banner by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}
