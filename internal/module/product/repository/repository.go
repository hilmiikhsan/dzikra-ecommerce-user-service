package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/ports"
	productCategoryDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	productGroceryDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/dto"
	productGrocery "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/entity"
	productImageDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/dto"
	productImage "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/entity"
	productSubCategoryDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/dto"
	productVariantDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/dto"
	productVariant "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/entity"
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

func (r *productRepository) UpdateProduct(ctx context.Context, tx *sqlx.Tx, id int, data *entity.Product) (*entity.Product, error) {
	var res = new(entity.Product)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateProduct),
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
		id,
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
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateProduct - product with id %d not found", id)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrProductNotFound)
		}

		log.Error().Err(err).Msg("repository::UpdateProduct - error updating product")
		return nil, err
	}

	return res, nil
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

func (r *productRepository) FindListProduct(ctx context.Context, limit, offset int, search string, categoryID, subcategoryID int) ([]dto.GetListProduct, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountListProduct), search, categoryID, subcategoryID); err != nil {
		log.Error().Err(err).Msg("repository::FindListProduct - error counting products")
		return nil, 0, fmt.Errorf("error counting products: %w", err)
	}

	var rows []entity.Product
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(queryFindListProduct), search, categoryID, subcategoryID, limit, offset); err != nil {
		log.Error().Err(err).Msg("repository::FindListProduct - error selecting list of products")
		return nil, 0, fmt.Errorf("error selecting list of products: %w", err)
	}

	productMap := make(map[int]*dto.GetListProduct, len(rows))
	ordered := make([]*dto.GetListProduct, 0, len(rows))

	for _, row := range rows {
		prod, exists := productMap[row.ID]
		if !exists {
			prod = &dto.GetListProduct{
				ID:            row.ID,
				Name:          row.Name,
				Description:   row.Description,
				Specification: row.Spesification,
				RealPrice:     row.RealPrice,
				CapitalPrice:  row.CapitalPrice,
				DiscountPrice: row.DiscountPrice,
				Stock:         row.Stock,
				Weight:        row.Weight,
				VariantName:   row.VariantName,
				ProductCategory: productCategoryDto.GetListCategory{
					ID:       row.ProductCategoryID,
					Category: row.ProductCategoryName,
				},
				ProductSubCategory: productSubCategoryDto.ProductSubCategory{
					ID:          row.ProductSubID,
					SubCategory: row.ProductSubCategoryName,
					CategoryID:  row.ProductSubCategoryID,
				},
				ProductVariant: []productVariantDto.ProductVariant{},
				ProductGrocery: []productGroceryDto.ProductGrocery{},
				ProductImage:   []productImageDto.ProductImage{},
			}

			productMap[row.ID] = prod
			ordered = append(ordered, prod)
		}

		if row.ProductVariantID.Valid {
			variantID := int(row.ProductVariantID.Int64)
			alreadyExists := false
			for _, v := range prod.ProductVariant {
				if v.ID == variantID {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				variant := productVariantDto.ProductVariant{
					ID:             variantID,
					VariantSubName: row.ProductVariantSubName.String,
					VariantStock:   int(row.ProductVariantStock.Int64),
					VariantWeight:  row.ProductVariantWeight.Float64,
					CapitalPrice:   int(row.ProductVariantCapitalPrice.Int64),
					RealPrice:      int(row.ProductVariantRealPrice.Int64),
					DiscountPrice:  int(row.ProductVariantDiscountPrice.Int64),
					ProductID:      int(row.ProductVariantProductID.Int64),
				}

				prod.ProductVariant = append(prod.ProductVariant, variant)
			}
		}

		if row.ProductGroceryID.Valid {
			groceryID := int(row.ProductGroceryID.Int64)
			alreadyExists := false
			for _, g := range prod.ProductGrocery {
				if g.ID == groceryID {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				grocery := productGroceryDto.ProductGrocery{
					ID:        groceryID,
					MinBuy:    int(row.ProductGroceryMinBuy.Int64),
					Discount:  int(row.ProductGroceryDiscount.Int64),
					ProductID: int(row.ProductGroceryProductID.Int64),
				}

				prod.ProductGrocery = append(prod.ProductGrocery, grocery)
			}
		}

		if row.ProductImageID.Valid {
			imageID := int(row.ProductImageID.Int64)
			alreadyExists := false
			for _, img := range prod.ProductImage {
				if img.ID == imageID {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				image := productImageDto.ProductImage{
					ID:        imageID,
					ImageURL:  row.ProductImageURL.String,
					Position:  int(row.ProductImageSort.Int64),
					ProductID: int(row.ProductImageProductID.Int64),
				}

				prod.ProductImage = append(prod.ProductImage, image)
			}
		}
	}

	results := make([]dto.GetListProduct, len(ordered))
	for i, p := range ordered {
		results[i] = *p
	}

	return results, total, nil
}

func (r *productRepository) FindProductByID(ctx context.Context, id int) (*entity.Product, error) {
	res := new(entity.Product)
	if err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindProductByID), id); err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Int("id", id).Msg("repository::FindProductByID - Failed to find product by id")
			return nil, errors.New(constants.ErrProductNotFound)
		}

		log.Error().Err(err).Msg("repository::FindProductByID - error fetching product")
		return nil, err
	}

	var variants []productVariant.ProductVariant
	if err := r.db.SelectContext(ctx, &variants, r.db.Rebind(queryFindVariants), id); err != nil && err != sql.ErrNoRows {
		log.Error().Err(err).Msg("repository::FindProductByID - error fetching product variants")
		return nil, err
	}
	res.ProductVariant = variants

	var groceries []productGrocery.ProductGrocery
	if err := r.db.SelectContext(ctx, &groceries, r.db.Rebind(queryFindGroceries), id); err != nil && err != sql.ErrNoRows {
		log.Error().Err(err).Msg("repository::FindProductByID - error fetching product groceries")
		return nil, err
	}
	res.ProductGrocery = groceries

	var images []productImage.ProductImage
	if err := r.db.SelectContext(ctx, &images, r.db.Rebind(queryFindImages), id); err != nil && err != sql.ErrNoRows {
		log.Error().Err(err).Msg("repository::FindProductByID - error fetching product images")
		return nil, err
	}
	res.ProductImage = images

	return res, nil
}

func (r *productRepository) SoftDeleteProductByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteProductByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::SoftDeleteProductByID - Failed to soft delete product")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteProductByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrProductNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::SoftDeleteProductByID - Product not found")
		return errNotFound
	}

	return nil
}
