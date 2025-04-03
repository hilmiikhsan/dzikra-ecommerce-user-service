package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductSubCategoryRepository = &productSubCategoryRepository{}

type productSubCategoryRepository struct {
	db *sqlx.DB
}

func NewProductSubCategoryRepository(db *sqlx.DB) *productSubCategoryRepository {
	return &productSubCategoryRepository{
		db: db,
	}
}

func (r *productSubCategoryRepository) InsertNewProductSubCategory(ctx context.Context, name string, categoryID int) (*entity.ProductSubCategory, error) {
	var res = new(entity.ProductSubCategory)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryInsertNewProductSubCategory), name, categoryID).Scan(&res.ID, &res.Name, &res.ProductCategoryID)
	if err != nil {
		uniqueConstraints := map[string]string{
			"product_sub_categories_name_key": constants.ErrProductSubCategoryAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, name, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", name).Msg("repository::InsertNewProductSubCategory - Failed to insert new product sub category")
			return nil, handleErr
		}

		if productSubCategory, ok := val.(*entity.ProductSubCategory); ok {
			log.Error().Err(err).Any("payload", name).Msg("repository::InsertNewProductSubCategory - Failed to insert new product sub category")
			return productSubCategory, nil
		}

		log.Error().Err(err).Str("name", name).Msg("repository::InsertNewProductSubCategory - error inserting new product sub category")
		return nil, err
	}

	return res, nil
}
