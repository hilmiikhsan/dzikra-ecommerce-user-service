package service

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *productVariantService) GetProductVariantStock(ctx context.Context, id int) (int, error) {
	stock, err := s.productVariantRepository.FindProductVariantStockByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("service::GetProductVariantStock - Failed to get product variant stock")
		return 0, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return stock, nil
}
