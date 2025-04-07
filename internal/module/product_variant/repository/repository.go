package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductVariantRepository = &productVariantRepository{}

type productVariantRepository struct {
	db *sqlx.DB
}

func NewProductVariantRepository(db *sqlx.DB) *productVariantRepository {
	return &productVariantRepository{
		db: db,
	}
}

func (r *productVariantRepository) InsertNewProductVariant(ctx context.Context, tx *sqlx.Tx, data *entity.ProductVariant) (*entity.ProductVariant, error) {
	var res = new(entity.ProductVariant)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewProductVariant),
		data.VariantSubName,
		data.VariantStock,
		data.VariantWeight,
		data.CapitalPrice,
		data.RealPrice,
		data.DiscountPrice,
		data.ProductID,
	).Scan(
		&res.ID,
		&res.VariantSubName,
		&res.VariantStock,
		&res.VariantWeight,
		&res.CapitalPrice,
		&res.RealPrice,
		&res.DiscountPrice,
		&res.ProductID,
	)
	if err != nil {
		log.Error().Err(err).Msg("repository::InsertNewProductVariant - error inserting new product variant")
		return nil, err
	}

	return res, nil
}
