package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/order/dto"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest, locals *middleware.Locals, addressID, voucherID int) (*dto.CreateOrderResponse, error)
	GetListOrder(ctx context.Context, page, limit int, search, status, userID string) (*dto.GetListOrderResponse, error)
	GetWaybillDetails(ctx context.Context, orderID string) (*dto.GetWaybillResponse, error)
	GetListOrderTransaction(ctx context.Context, page, limit int, search, status string) (*dto.GetListOrderResponse, error)
	UpdateOrderShippingNumber(ctx context.Context, req *dto.UpdateOrderShippingNumberRequest, orderID string) (*dto.UpdateOrderShippingNumberResponse, error)
	UpdateOrderStatusTransaction(ctx context.Context, req *dto.UpdateOrderStatusTransactionRequest, orderID, fullName, email string) error
}
