package service

import (
	productCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/ports"
	productSubCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/ports"
	"github.com/jmoiron/sqlx"
)

var _ productSubCategoryPorts.ProductSubCategoryService = &productSubCategoryService{}

type productSubCategoryService struct {
	db                           *sqlx.DB
	productSubCategoryRepository productSubCategoryPorts.ProductSubCategoryRepository
	productCategoryRepository    productCategoryPorts.ProductCategoryRepository
}

func NewProductSubCategoryService(
	db *sqlx.DB,
	productSubCategoryRepository productSubCategoryPorts.ProductSubCategoryRepository,
	productCategoryRepository productCategoryPorts.ProductCategoryRepository,
) *productSubCategoryService {
	return &productSubCategoryService{
		db:                           db,
		productSubCategoryRepository: productSubCategoryRepository,
		productCategoryRepository:    productCategoryRepository,
	}
}
