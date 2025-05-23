package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/entity"
	"github.com/jmoiron/sqlx"
)

type ProductGroceryRepository interface {
	InsertNewProductGrocery(ctx context.Context, tx *sqlx.Tx, data *entity.ProductGrocery) (*entity.ProductGrocery, error)
	UpdateProductGrocery(ctx context.Context, tx *sqlx.Tx, data *entity.ProductGrocery) (*entity.ProductGrocery, error)
	SoftDeleteProductGroceriesByProductID(ctx context.Context, tx *sqlx.Tx, productID int) error
	FindProductGroceryByProductID(ctx context.Context, productID int) ([]dto.GroceryPrice, error)
}
