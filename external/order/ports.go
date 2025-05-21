package order

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/proto/order"
)

type ExternalOrder interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
	GetListOrder(ctx context.Context, page, limit int, search, status, userID string) (*order.GetListOrderResponse, error)
	GetOrderById(ctx context.Context, req *order.GetOrderByIdRequest) (*order.GetOrderByIdResponse, error)
	GetListOrderTransaction(ctx context.Context, page, limit int, search, status string) (*order.GetListOrderResponse, error)
	UpdateOrderShippingNumber(ctx context.Context, req *order.UpdateOrderShippingNumberRequest) (*order.UpdateOrderShippingNumberResponse, error)
	UpdateOrderStatusTransaction(ctx context.Context, req *order.UpdateOrderStatusTransactionRequest) (*order.UpdateOrderStatusTransactionResponse, error)
	GetOrderItemsByOrderID(ctx context.Context, req *order.GetOrderItemsByOrderIDRequest) (*order.GetOrderItemsByOrderIDResponse, error)
}
