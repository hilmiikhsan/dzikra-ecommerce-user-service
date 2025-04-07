package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/entity"
	"github.com/jmoiron/sqlx"
)

type ProductVariantRepository interface {
	InsertNewProductVariant(ctx context.Context, tx *sqlx.Tx, data *entity.ProductVariant) (*entity.ProductVariant, error)
}
