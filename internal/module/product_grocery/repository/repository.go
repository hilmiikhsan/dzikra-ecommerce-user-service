package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductGroceryRepository = &productGroceryRepository{}

type productGroceryRepository struct {
	db *sqlx.DB
}

func NewProductGroceryRepository(db *sqlx.DB) *productGroceryRepository {
	return &productGroceryRepository{
		db: db,
	}
}

func (r *productGroceryRepository) InsertNewProductGrocery(ctx context.Context, tx *sqlx.Tx, data *entity.ProductGrocery) (*entity.ProductGrocery, error) {
	var res = new(entity.ProductGrocery)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewProductGrocery),
		data.MinBuy,
		data.Discount,
		data.ProductID,
	).Scan(
		&res.ID,
		&res.MinBuy,
		&res.Discount,
		&res.ProductID,
	)
	if err != nil {
		log.Error().Err(err).Msg("repository::InsertNewProductGrocery - error inserting new product grocery")
		return nil, err
	}

	return res, nil
}

func (r *productGroceryRepository) UpdateProductGrocery(ctx context.Context, tx *sqlx.Tx, data *entity.ProductGrocery) (*entity.ProductGrocery, error) {
	var res = new(entity.ProductGrocery)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateProductGrocery),
		data.MinBuy,
		data.Discount,
		data.ID,
		data.ProductID,
	).Scan(
		&res.ID,
		&res.MinBuy,
		&res.Discount,
		&res.ProductID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("repository::UpdateProductGrocery - product grocery not found or no update performed")
			return nil, fmt.Errorf("product grocery with id %d and product_id %d not found", data.ID, data.ProductID)
		}

		log.Error().Err(err).Msg("repository::UpdateProductGrocery - error updating product grocery")
		return nil, err
	}

	return res, nil
}

func (r *productGroceryRepository) SoftDeleteProductGroceriesByProductID(ctx context.Context, tx *sqlx.Tx, productID int) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteProductGroceriesByProductID), productID)
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteProductGroceriesByProductID - error soft deleting product groceries")
		return err
	}

	return nil
}
