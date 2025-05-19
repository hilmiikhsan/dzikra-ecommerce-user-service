package grpc

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/cart"
	"github.com/rs/zerolog/log"
)

func (h *cartGrpcAPI) DeleteCartByUserID(ctx context.Context, req *cart.DeleteCartByUserIdRequest) (*cart.DeleteCartByUserIdResponse, error) {
	err := h.CartService.DeleteCartItemByUserID(ctx, req.UserId)
	if err != nil {
		log.Err(err).Msg("order::DeleteCartByUserID - Failed to delete cart item")
		return &cart.DeleteCartByUserIdResponse{
			Message: "failed to delete cart item",
		}, nil
	}

	return &cart.DeleteCartByUserIdResponse{
		Message: "success",
	}, nil
}
