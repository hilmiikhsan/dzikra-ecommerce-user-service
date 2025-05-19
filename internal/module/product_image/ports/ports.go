package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/entity"
	"github.com/jmoiron/sqlx"
)

type ProductImageRepository interface {
	InsertNewProductImage(ctx context.Context, tx *sqlx.Tx, data *entity.ProductImage) (*entity.ProductImage, error)
	GetNextSort(ctx context.Context, productID int) (int, error)
	ReorderProductImages(ctx context.Context, tx *sqlx.Tx, productID int) error
	CountProductImagesByProductID(ctx context.Context, productID int) (int, error)
	UpdateProductImageURL(ctx context.Context, id int, url string) (*entity.ProductImage, error)
	FindProductImagesByProductID(ctx context.Context, productID int) ([]entity.ProductImage, error)
	DeleteProductImage(ctx context.Context, tx *sqlx.Tx, id int) error
	SoftDeleteProductImagesByProductID(ctx context.Context, tx *sqlx.Tx, productID int) error
	FindImagesByProductIds(ctx context.Context, productIDs []int64) ([]entity.ProductImage, error)
}

type ProductImageService interface {
	GetImagesByProductIds(ctx context.Context, productIDs []int64) ([]dto.ProductImage, error)
}
