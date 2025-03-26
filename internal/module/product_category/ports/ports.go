package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
)

type ProductCategoryRepository interface {
	FindListProductCategory(ctx context.Context, limit, offset int, search string) ([]dto.GetListCategory, int, error)
}

type ProductCategoryService interface {
	GetListProductCategory(ctx context.Context, page, limit int, search string) (*dto.GetListProductCategory, error)
}
