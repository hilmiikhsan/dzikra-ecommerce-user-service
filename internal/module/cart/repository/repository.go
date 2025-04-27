package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.CartRepository = &cartRepository{}

type cartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) *cartRepository {
	return &cartRepository{
		db: db,
	}
}

func (r *cartRepository) InsertNewCart(ctx context.Context, tx *sqlx.Tx, data *entity.Cart) (*entity.Cart, error) {
	var res = new(entity.Cart)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewCart),
		data.UserID,
		data.ProductID,
		data.ProductVariantID,
		data.Quantity,
	).Scan(
		&res.ID,
		&res.UserID,
		&res.ProductID,
		&res.ProductVariantID,
		&res.Quantity,
		&res.CreatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewCart - Failed to insert new cart")
		return nil, err
	}

	return res, nil
}
