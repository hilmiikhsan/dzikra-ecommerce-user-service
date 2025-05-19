package service

import (
	productImagePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/ports"
	"github.com/jmoiron/sqlx"
)

var _ productImagePorts.ProductImageService = &productImageService{}

type productImageService struct {
	db                     *sqlx.DB
	productImageRepository productImagePorts.ProductImageRepository
}

func NewProductImageService(
	db *sqlx.DB,
	productImageRepository productImagePorts.ProductImageRepository,
) *productImageService {
	return &productImageService{
		db:                     db,
		productImageRepository: productImageRepository,
	}
}
