package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/entity"
	"github.com/jmoiron/sqlx"
)

type ProductImageRepository interface {
	InsertNewProductImage(ctx context.Context, tx *sqlx.Tx, data *entity.ProductImage) (*entity.ProductImage, error)
	GetNextSort(ctx context.Context, productID int) (int, error)
	ReorderProductImages(ctx context.Context, tx *sqlx.Tx, productID int) error
	CountProductImagesByProductID(ctx context.Context, productID int) (int, error)
	UpdateProductImageURL(ctx context.Context, id int, url string) (*entity.ProductImage, error)
}
