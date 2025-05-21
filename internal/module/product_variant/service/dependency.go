package service

import (
	productVariantPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/ports"
)

var _ productVariantPorts.ProductVariantService = &productVariantService{}

type productVariantService struct {
	productVariantRepository productVariantPorts.ProductVariantRepository
}

func NewProductVariantService(
	productVariantRepository productVariantPorts.ProductVariantRepository,
) *productVariantService {
	return &productVariantService{
		productVariantRepository: productVariantRepository,
	}
}
