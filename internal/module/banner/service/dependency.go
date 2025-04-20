package service

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/storage/minio"
	bannerPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/ports"
)

var _ bannerPorts.BannerService = &bannerService{}

type bannerService struct {
	bannerRepository bannerPorts.BannerRepository
	minioService     minio.MinioService
}

func NewBannerService(
	bannerRepository bannerPorts.BannerRepository,
	minioService minio.MinioService,
) *bannerService {
	return &bannerService{
		bannerRepository: bannerRepository,
		minioService:     minioService,
	}
}
