package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/entity"
	"github.com/jmoiron/sqlx"
)

type ProductGroceryRepository interface {
	InsertNewProductGrocery(ctx context.Context, tx *sqlx.Tx, data *entity.ProductGrocery) (*entity.ProductGrocery, error)
}
