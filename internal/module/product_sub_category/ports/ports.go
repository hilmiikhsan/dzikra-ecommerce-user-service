package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/entity"
)

type ProductSubCategoryRepository interface {
	InsertNewProductSubCategory(ctx context.Context, name string, categoryID int) (*entity.ProductSubCategory, error)
}

type ProductSubCategoryService interface {
	CreateProductSubCategory(ctx context.Context, req *dto.CreateOrUpdateProductSubCategoryRequest, categoryID int) (*dto.CreateOrUpdateProductSubCategoryResponse, error)
}
