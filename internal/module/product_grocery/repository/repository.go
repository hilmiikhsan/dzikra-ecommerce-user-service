package repository

import (
	"context"

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
