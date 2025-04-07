package service

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/storage/minio"
	productPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/ports"
	productCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/ports"
	productGroceryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/ports"
	productImagePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/ports"
	productSubCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/ports"
	productVariantPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/ports"
	"github.com/jmoiron/sqlx"
)

var _ productPorts.ProductService = &productService{}

type productService struct {
	db                           *sqlx.DB
	productRepository            productPorts.ProductRepository
	productCategoryRepository    productCategoryPorts.ProductCategoryRepository
	productSubCategoryRepository productSubCategoryPorts.ProductSubCategoryRepository
	productVariantRepository     productVariantPorts.ProductVariantRepository
	productGroceryRepository     productGroceryPorts.ProductGroceryRepository
	minioService                 minio.MinioService
	productImageRepository       productImagePorts.ProductImageRepository
}

func NewProductService(
	db *sqlx.DB,
	productRepository productPorts.ProductRepository,
	productCategoryRepository productCategoryPorts.ProductCategoryRepository,
	productSubCategoryRepository productSubCategoryPorts.ProductSubCategoryRepository,
	productVariantRepository productVariantPorts.ProductVariantRepository,
	productGroceryRepository productGroceryPorts.ProductGroceryRepository,
	minioService minio.MinioService,
	productImageRepository productImagePorts.ProductImageRepository,
) *productService {
	return &productService{
		db:                           db,
		productRepository:            productRepository,
		productCategoryRepository:    productCategoryRepository,
		productSubCategoryRepository: productSubCategoryRepository,
		productVariantRepository:     productVariantRepository,
		productGroceryRepository:     productGroceryRepository,
		minioService:                 minioService,
		productImageRepository:       productImageRepository,
	}
}
