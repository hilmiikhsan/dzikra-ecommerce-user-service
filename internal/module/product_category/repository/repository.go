package repostiory

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductCategoryRepository = &productCategoryRepository{}

type productCategoryRepository struct {
	db *sqlx.DB
}

func NewProductCategoryRepository(db *sqlx.DB) *productCategoryRepository {
	return &productCategoryRepository{
		db: db,
	}
}

func (r *productCategoryRepository) FindListProductCategory(ctx context.Context, limit, offset int, search string) ([]dto.GetListCategory, int, error) {
	var responses []entity.ProductCategory

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListProductCategory), search, limit, offset); err != nil {
		log.Error().Err(err).Msg("repository::FindListProductCategory - error executing query")
		return nil, 0, err
	}

	var total int

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountListProductCategory), search); err != nil {
		log.Error().Err(err).Msg("repository::FindListProductCategory - error counting product categories")
		return nil, 0, err
	}

	var productCategories []dto.GetListCategory

	for _, response := range responses {
		productCategoryDTO := dto.GetListCategory{
			ID:       response.ID,
			Category: response.Name,
		}

		productCategories = append(productCategories, productCategoryDTO)
	}

	return productCategories, total, nil
}

func (r *productCategoryRepository) InsertNewProductCategory(ctx context.Context, name string) (*entity.ProductCategory, error) {
	var res = new(entity.ProductCategory)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryInsertNewProductCategory), name).Scan(&res.ID, &res.Name)
	if err != nil {
		uniqueConstraints := map[string]string{
			"product_categories_name_key": constants.ErrProductCategoryAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, name, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", name).Msg("repository::InsertNewProductCategory - Failed to insert new product category")
			return nil, handleErr
		}

		if productCategory, ok := val.(*entity.ProductCategory); ok {
			log.Error().Err(err).Any("payload", name).Msg("repository::InsertNewProductCategory - Failed to insert new product category")
			return productCategory, nil
		}

		log.Error().Err(err).Str("name", name).Msg("repository::InsertNewProductCategory - error inserting new product category")
		return nil, err
	}

	return res, nil
}

func (r *productCategoryRepository) UpdateProductCategory(ctx context.Context, id int, name string) (*entity.ProductCategory, error) {
	var res = new(entity.ProductCategory)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryUpdateProductCategory), name, id, name).Scan(&res.ID, &res.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Int("id", id).Str("name", name).Msg("repository::UpdateProductCategory - product category is already registered")
			return nil, errors.New(constants.ErrProductCategoryAlreadyRegistered)
		}

		log.Error().Err(err).Str("name", name).Msg("repository::UpdateProductCategory - error updating product category")
		return nil, err
	}

	return res, nil
}

func (r *productCategoryRepository) FindProductCategoryByID(ctx context.Context, id int) (*entity.ProductCategory, error) {
	var res = new(entity.ProductCategory)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindProductCategoryByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Int("id", id).Msg("repository::FindProductCategoryByID - Failed to find product category by id")
			return nil, errors.New(constants.ErrProductCategoryNotFound)
		}

		log.Error().Err(err).Int("id", id).Msg("repository::FindProductCategoryByID - error finding product category by id")
		return nil, err
	}

	return res, nil
}

func (r *productCategoryRepository) DeleteProductCategoryByID(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, r.db.Rebind(queryDeleteProductCategoryByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::DeleteProductCategoryByID - Failed to soft delete Product category")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::DeleteProductCategoryByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrProductCategoryNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::DeleteProductCategoryByID - Product Category not found")
		return errNotFound
	}

	return nil
}

func (r *productCategoryRepository) CountProductCategoryByID(ctx context.Context, id int) (int, error) {
	var count int

	err := r.db.GetContext(ctx, &count, r.db.Rebind(queryCountProductCategoryByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::CountProductCategoryByID - error counting product category by id")
		return 0, err
	}

	return count, nil
}
