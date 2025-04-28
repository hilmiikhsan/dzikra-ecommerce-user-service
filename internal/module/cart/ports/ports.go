package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CartRepository interface {
	InsertNewCart(ctx context.Context, tx *sqlx.Tx, data *entity.Cart) (*entity.Cart, error)
	FindListCartByUserID(ctx context.Context, userID uuid.UUID) ([]dto.GetListCartResponse, error)
}

type CartService interface {
	AddCartItem(ctx context.Context, req *dto.AddCartItemRequest) (*dto.AddCartItemResponse, error)
	GetListCart(ctx context.Context, userID string) (*[]dto.GetListCartResponse, error)
}
