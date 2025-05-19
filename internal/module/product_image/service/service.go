package service

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *productImageService) GetImagesByProductIds(ctx context.Context, productIDs []int64) ([]dto.ProductImage, error) {
	images, err := s.productImageRepository.FindImagesByProductIds(ctx, productIDs)
	if err != nil {
		log.Error().Err(err).Msg("service::GetImagesByProductIds - error fetching images by product IDs")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var result []dto.ProductImage
	for _, image := range images {
		result = append(result, dto.ProductImage{
			ID:        image.ID,
			ImageURL:  image.ImageURL,
			Position:  image.Sort,
			ProductID: image.ProductID,
		})
	}

	return result, nil
}
