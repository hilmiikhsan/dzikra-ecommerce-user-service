package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/entity"
)

type BannerRepository interface {
	InsertNewBanner(ctx context.Context, data *entity.Banner) (*entity.Banner, error)
	FindListBanner(ctx context.Context, limit, offset int, search string) ([]dto.GetListBanner, int, error)
}

type BannerService interface {
	CreateBanner(ctx context.Context, description string, payloadFile dto.UploadFileRequest) (*dto.CreateBannerResponse, error)
	GetListBanner(ctx context.Context, page, limit int, search string) (*dto.GetListBannerResponse, error)
}
