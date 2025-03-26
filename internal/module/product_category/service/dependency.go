package service

import (
	productCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/ports"
	"github.com/jmoiron/sqlx"
)

var _ productCategoryPorts.ProductCategoryService = &productCategoryService{}

type productCategoryService struct {
	db                        *sqlx.DB
	productCategoryRepository productCategoryPorts.ProductCategoryRepository
}

func NewProductCategoryService(
	db *sqlx.DB,
	productCategoryRepository productCategoryPorts.ProductCategoryRepository,
) *productCategoryService {
	return &productCategoryService{
		db:                        db,
		productCategoryRepository: productCategoryRepository,
	}
}
