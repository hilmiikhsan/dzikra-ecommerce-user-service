package order

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/proto/order"
)

type ExternalOrder interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
	GetListOrder(ctx context.Context, page, limit int, search, status, userID string) (*order.GetListOrderResponse, error)
}
