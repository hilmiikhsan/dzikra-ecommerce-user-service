package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/entity"
)

type BannerRepository interface {
	InsertNewBanner(ctx context.Context, data *entity.Banner) (*entity.Banner, error)
}

type BannerService interface {
	CreateBanner(ctx context.Context, description string, payloadFile dto.UploadFileRequest) (*dto.CreateBannerResponse, error)
}
