package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	productCategory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/dto"
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

func (r *productSubCategoryRepository) FindProductSubCategoryByID(ctx context.Context, id int) (*entity.ProductSubCategory, error) {
	var res = new(entity.ProductSubCategory)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindProductSubCategoryByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Int("id", id).Msg("repository::FindProductSubCategoryByProductCategoryID - Failed to find product sub category by product category id")
			return nil, errors.New(constants.ErrProductSubCategoryNotFound)
		}

		log.Error().Err(err).Int("id", id).Msg("repository::FindProductSubCategoryByProductCategoryID - error finding product sub category by product category id")
		return nil, err
	}

	return res, nil
}

func (r *productSubCategoryRepository) UpdateProducSubCategory(ctx context.Context, name string, subCategoryID int) (*entity.ProductSubCategory, error) {
	var res = new(entity.ProductSubCategory)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryUpdateProductSubCategory), name, subCategoryID, name).Scan(&res.ID, &res.Name, &res.ProductCategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Int("product_category_id", subCategoryID).Str("name", name).Msg("repository::UpdateProducSubCategory - product sub category is already registered")
			return nil, errors.New(constants.ErrProductSubCategoryAlreadyRegistered)
		}

		log.Error().Err(err).Str("name", name).Msg("repository::UpdateProducSubCategory - error updating product sub category")
		return nil, err
	}

	return res, nil
}

func (r *productSubCategoryRepository) FindListProductSubCategory(ctx context.Context, limit, offset, categoryID int, search string) ([]dto.GetListSubCategory, int, error) {
	var responses []entity.ProductSubCategory

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListProductSubCategory), categoryID, search, limit, offset); err != nil {
		utils.QueryLog(r.db.Rebind(queryFindListProductSubCategory), limit, offset, categoryID, search)
		log.Error().Err(err).Msg("repository::FindListProductSubCategory - error finding list product sub category")
		return nil, 0, err
	}

	var total int

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountListProductSubCategory), search); err != nil {
		log.Error().Err(err).Msg("repository::FindListProductSubCategory - error counting list product sub category")
		return nil, 0, err
	}

	var productSubCategories []dto.GetListSubCategory

	for _, response := range responses {
		productSubCategoryDTO := dto.GetListSubCategory{
			ID:          response.ID,
			SubCategory: response.Name,
			Category: productCategory.GetListCategory{
				ID:       response.CategoryID,
				Category: response.CategoryName,
			},
		}

		productSubCategories = append(productSubCategories, productSubCategoryDTO)
	}

	return productSubCategories, total, nil
}

func (r *productSubCategoryRepository) SoftDeleteProductSubCategory(ctx context.Context, subCategoryID int) error {
	_, err := r.db.ExecContext(ctx, r.db.Rebind(querySoftDeleteProductSubCategory), subCategoryID)
	if err != nil {
		log.Error().Err(err).Int("product_sub_category_id", subCategoryID).Msg("repository::DeleteProductSubCategory - error deleting product sub category")
		return err
	}

	return nil
}
