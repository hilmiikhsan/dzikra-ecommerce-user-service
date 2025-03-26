package repostiory

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/ports"
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

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListRole), search, limit, offset); err != nil {
		log.Error().Err(err).Msg("repository::FindListProductCategory - error executing query")
		return nil, 0, err
	}

	var total int

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountListRole), search); err != nil {
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
