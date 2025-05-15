package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/shipping/dto"
)

type ShippingService interface {
	CalculateShippingCost(ctx context.Context, req *dto.CalculateShippingCostRequest, userID string) ([]dto.CalculateShippingCostResponse, error)
}
