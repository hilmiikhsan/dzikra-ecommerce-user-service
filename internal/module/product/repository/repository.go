package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepository{}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CountProductByName(ctx context.Context, name string) (int, error) {
	var count int

	err := r.db.GetContext(ctx, &count, queryCountProductByName, name)
	if err != nil {
		log.Error().Err(err).Msg("repository::CountProductByName - error count product query")
		return 0, err
	}

	return count, nil
}

func (r *productRepository) InsertNewProduct(ctx context.Context, tx *sqlx.Tx, data *entity.Product) (*entity.Product, error) {
	var res = new(entity.Product)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewProduct),
		data.Name,
		data.RealPrice,
		data.DiscountPrice,
		data.CapitalPrice,
		data.Description,
		data.Spesification,
		data.Stock,
		data.Weight,
		data.VariantName,
		data.ProductCategoryID,
		data.ProductSubCategoryID,
	).Scan(
		&res.ID,
		&res.Name,
		&res.RealPrice,
		&res.DiscountPrice,
		&res.CapitalPrice,
		&res.Description,
		&res.Spesification,
		&res.Stock,
		&res.Weight,
		&res.VariantName,
		&res.ProductCategoryID,
		&res.ProductSubCategoryID,
	)
	if err != nil {
		log.Error().Err(err).Msg("repository::InsertNewProduct - error inserting new product")
		return nil, err
	}

	return res, nil
}
