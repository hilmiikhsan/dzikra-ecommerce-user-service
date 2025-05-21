package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/entity"
	"github.com/jmoiron/sqlx"
)

type ProductVariantRepository interface {
	InsertNewProductVariant(ctx context.Context, tx *sqlx.Tx, data *entity.ProductVariant) (*entity.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, tx *sqlx.Tx, data *entity.ProductVariant) (*entity.ProductVariant, error)
	DeleteProductVariant(ctx context.Context, tx *sqlx.Tx, id, productID int) error
	SoftDeleteProductVariantsByProductID(ctx context.Context, tx *sqlx.Tx, productID int) error
	CountProductVariantByIDAndProductID(ctx context.Context, id, productID int) (int, error)
	FindProductVariantStockByID(ctx context.Context, id int) (int, error)
}

type ProductVariantService interface {
	GetProductVariantStock(ctx context.Context, id int) (int, error)
}
