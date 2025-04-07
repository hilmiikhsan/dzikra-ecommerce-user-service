package service

import (
	productCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/ports"
)

var _ productCategoryPorts.ProductCategoryService = &productCategoryService{}

type productCategoryService struct {
	productCategoryRepository productCategoryPorts.ProductCategoryRepository
}

func NewProductCategoryService(
	productCategoryRepository productCategoryPorts.ProductCategoryRepository,
) *productCategoryService {
	return &productCategoryService{
		productCategoryRepository: productCategoryRepository,
	}
}
