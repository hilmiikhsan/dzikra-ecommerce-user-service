package grpc

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/product_variant"
	"github.com/rs/zerolog/log"
)

func (h *productVariantGrpcAPI) GetProductVariantStock(ctx context.Context, req *product_variant.GetProductVariantStockRequest) (*product_variant.GetProductVariantStockResponse, error) {
	stock, err := h.ProductVariantService.GetProductVariantStock(ctx, int(req.Id))
	if err != nil {
		log.Err(err).Msg("product_variant::GetProductVariantStock - Failed to get product variant stock")
		return &product_variant.GetProductVariantStockResponse{
			Message: "failed to get product variant stock",
		}, nil
	}

	return &product_variant.GetProductVariantStockResponse{
		Stock:   int64(stock),
		Message: "success",
	}, nil
}
