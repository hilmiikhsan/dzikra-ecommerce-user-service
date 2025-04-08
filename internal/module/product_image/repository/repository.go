package repository

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductImageRepository = &productImageRepository{}

type productImageRepository struct {
	db *sqlx.DB
}

func NewProductImageRepository(db *sqlx.DB) *productImageRepository {
	return &productImageRepository{
		db: db,
	}
}

func (r *productImageRepository) InsertNewProductImage(ctx context.Context, tx *sqlx.Tx, data *entity.ProductImage) (*entity.ProductImage, error) {
	if data.Sort == 0 {
		nextSort, err := r.GetNextSort(ctx, data.ProductID)
		if err != nil {
			log.Error().Err(err).Int("productID", data.ProductID).Msg("repository::InsertNewProductImage - error fetching next sort value")
			return nil, err
		}

		data.Sort = nextSort
	}

	var res = new(entity.ProductImage)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewProductImage),
		data.ImageURL,
		data.ProductID,
		data.Sort,
	).Scan(
		&res.ID,
		&res.ProductID,
		&res.Sort,
	)
	if err != nil {
		log.Error().Err(err).Msg("repository::InsertNewProductImage - error inserting new product image")
		return nil, err
	}

	return res, nil
}

func (r *productImageRepository) GetNextSort(ctx context.Context, productID int) (int, error) {
	var nextSort int

	err := r.db.GetContext(ctx, &nextSort, r.db.Rebind(queryGetNextSort), productID)
	if err != nil {
		log.Error().Err(err).Int("productID", productID).Msg("repository::GetNextSort - error fetching next sort value")
		return 0, err
	}

	return nextSort, nil
}

func (r *productImageRepository) ReorderProductImages(ctx context.Context, tx *sqlx.Tx, productID int) error {
	var images []entity.ProductImage

	err := tx.SelectContext(ctx, &images, r.db.Rebind(queryFindOrderProductImage), productID)
	if err != nil {
		log.Error().Err(err).Int("productID", productID).Msg("repository::ReorderProductImages - error fetching active images")
		return err
	}

	for index, img := range images {
		newSort := index + 1

		_, err := tx.ExecContext(ctx, r.db.Rebind(queryUpdateProductImageSorting), newSort, img.ID)
		if err != nil {
			log.Error().Err(err).Int("imgID", img.ID).Msg("repository::ReorderProductImages - error updating sort")
			return err
		}
	}

	return nil
}

func (r *productImageRepository) CountProductImagesByProductID(ctx context.Context, productID int) (int, error) {
	var count int

	if err := r.db.GetContext(ctx, &count, r.db.Rebind(queryCountProductImages), productID); err != nil {
		return 0, fmt.Errorf("repository::CountProductImagesByProductID - error counting product images: %w", err)
	}

	return count, nil
}

func (r *productImageRepository) UpdateProductImageURL(ctx context.Context, id int, url string) (*entity.ProductImage, error) {
	var res = new(entity.ProductImage)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryUpdateProductImageURL), url, id).Scan(
		&res.ID,
		&res.ImageURL,
		&res.ProductID,
		&res.Sort,
	)
	if err != nil {
		log.Error().Err(err).Msg("repository::UpdateProductImageURL - error updating image URL")
		return nil, err
	}

	return res, nil
}

func (r *productImageRepository) FindProductImagesByProductID(ctx context.Context, productID int) ([]entity.ProductImage, error) {
	var res []entity.ProductImage

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(queryFindProductImagesByProductID), productID)
	if err != nil {
		log.Error().Err(err).Msg("repository::FindProductImagesByProductID - error fetching product images")
		return nil, err
	}

	return res, nil
}

func (r *productImageRepository) DeleteProductImage(ctx context.Context, tx *sqlx.Tx, id int) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryDeleteProductImage), id)
	if err != nil {
		log.Error().Err(err).Msg("repository::DeleteProductImage - error deleting product image")
		return err
	}

	return nil
}
