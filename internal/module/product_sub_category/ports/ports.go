package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/entity"
)

type ProductSubCategoryRepository interface {
	InsertNewProductSubCategory(ctx context.Context, name string, categoryID int) (*entity.ProductSubCategory, error)
	FindProductSubCategoryByID(ctx context.Context, id int) (*entity.ProductSubCategory, error)
	UpdateProducSubCategory(ctx context.Context, name string, subCategoryID int) (*entity.ProductSubCategory, error)
	FindListProductSubCategory(ctx context.Context, limit, offset, categoryID int, search string) ([]dto.GetListSubCategory, int, error)
}

type ProductSubCategoryService interface {
	CreateProductSubCategory(ctx context.Context, req *dto.CreateOrUpdateProductSubCategoryRequest, categoryID int) (*dto.CreateOrUpdateProductSubCategoryResponse, error)
	UpdateProductSubCategory(ctx context.Context, req *dto.CreateOrUpdateProductSubCategoryRequest, categoryID, subCategoryID int) (*dto.CreateOrUpdateProductSubCategoryResponse, error)
	GetListProductSubCategory(ctx context.Context, page, limit, categoryID int, search string) (*dto.GetListProductSubCategory, error)
}
