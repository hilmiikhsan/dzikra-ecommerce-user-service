package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/entity"
)

type ProductCategoryRepository interface {
	FindListProductCategory(ctx context.Context, limit, offset int, search string) ([]dto.GetListCategory, int, error)
	InsertNewProductCategory(ctx context.Context, name string) (*entity.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, id int, name string) (*entity.ProductCategory, error)
	FindProductCategoryByID(ctx context.Context, id int) (*entity.ProductCategory, error)
	DeleteProductCategoryByID(ctx context.Context, id int) error
	CountProductCategoryByID(ctx context.Context, id int) (int, error)
}

type ProductCategoryService interface {
	GetListProductCategory(ctx context.Context, page, limit int, search string) (*dto.GetListProductCategoryResponse, error)
	GetDetailProductCategory(ctx context.Context, id int) (*dto.GetListCategory, error)
	CreateProductCategory(ctx context.Context, name string) (*dto.CreateOrProductCategoryResponse, error)
	UpdateProductCategory(ctx context.Context, id int, name string) (*dto.CreateOrProductCategoryResponse, error)
	RemoveProductCategory(ctx context.Context, id int) error
}
