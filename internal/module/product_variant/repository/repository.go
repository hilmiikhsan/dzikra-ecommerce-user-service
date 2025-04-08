package repository

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *productVariantRepository) UpdateProductVariant(ctx context.Context, tx *sqlx.Tx, data *entity.ProductVariant) (*entity.ProductVariant, error) {
	var res = new(entity.ProductVariant)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateProductVariant),
		data.VariantSubName,
		data.VariantStock,
		data.VariantWeight,
		data.CapitalPrice,
		data.RealPrice,
		data.DiscountPrice,
		data.ID,
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
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("repository::UpdateProductVariant - product variant not found or no update performed")
			return nil, fmt.Errorf("product variant with id %d and product_id %d not found", data.ID, data.ProductID)
		}

		log.Error().Err(err).Msg("repository::UpdateProductVariant - error updating product variant")
		return nil, err
	}

	return res, nil
}

func (r *productVariantRepository) DeleteProductVariant(ctx context.Context, tx *sqlx.Tx, id, productID int) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryDeleteProductVariant), id, productID)
	if err != nil {
		log.Error().Err(err).Msg("repository::DeleteProductVariant - error deleting product variant")
		return err
	}

	return nil
}
