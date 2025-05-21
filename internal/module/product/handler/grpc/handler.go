package grpc

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/product"
	"github.com/rs/zerolog/log"
)

func (h *productGrpcAPI) GetProductStock(ctx context.Context, req *product.GetProductStockRequest) (*product.GetProductStockResponse, error) {
	stock, err := h.ProductService.GetProductStock(ctx, int(req.Id))
	if err != nil {
		log.Err(err).Msg("product::GetProductStock - Failed to get product stock")
		return &product.GetProductStockResponse{
			Message: "failed to get product stock",
		}, nil
	}

	return &product.GetProductStockResponse{
		Stock:   int64(stock),
		Message: "success",
	}, nil
}
