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
	UpdateCart(ctx context.Context, tx *sqlx.Tx, data *entity.Cart) (*entity.Cart, error)
	DeleteCartByID(ctx context.Context, tx *sqlx.Tx, id int) error
}

type CartService interface {
	AddCartItem(ctx context.Context, req *dto.AddOrUpdateCartItemRequest) (*dto.AddOrUpdateCartItemResponse, error)
	GetListCartItem(ctx context.Context, userID string) (*[]dto.GetListCartResponse, error)
	UpdateCartItem(ctx context.Context, req *dto.AddOrUpdateCartItemRequest, id int) (*dto.AddOrUpdateCartItemResponse, error)
	DeleteCartItem(ctx context.Context, id int) error
}
