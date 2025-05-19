package grpc

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/product_image"
	"github.com/rs/zerolog/log"
)

func (h *productImageGrpcAPI) GetImagesByProductIds(ctx context.Context, req *product_image.GetImagesRequest) (*product_image.GetImagesResponse, error) {
	images, err := h.ProductImageService.GetImagesByProductIds(ctx, req.ProductIds)
	if err != nil {
		log.Err(err).Msg("product_image::GetImagesByProductIds - failed to get images by product IDs")
		return &product_image.GetImagesResponse{
			Message: "failed to get images by product IDs",
			Images:  nil,
		}, nil
	}

	resp := &product_image.GetImagesResponse{
		Message: "success",
	}

	for _, img := range images {
		resp.Images = append(resp.Images, &product_image.ProductImage{
			Id:        int64(img.ID),
			ImageUrl:  img.ImageURL,
			Position:  int32(img.Position),
			ProductId: int64(img.ProductID),
		})
	}

	return resp, nil
}
