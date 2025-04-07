package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/storage/minio"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/ports"
	productRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/repository"
	productService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/service"
	productCategoryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/repository"
	productGroceryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/repository"
	productImageRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/repository"
	productSubCategoryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/repository"
	productVariantRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/repository"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type productHandler struct {
	service    ports.ProductService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewProductHandler() *productHandler {
	var handler = new(productHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// minio service
	minioService := minio.NewMinioService(adapter.Adapters.DzikraMinio, config.Envs.MinioStorage.Bucket)

	// repository
	productRepository := productRepository.NewProductRepository(adapter.Adapters.DzikraPostgres)
	productCategoryRepository := productCategoryRepository.NewProductCategoryRepository(adapter.Adapters.DzikraPostgres)
	productSubCategoryRepository := productSubCategoryRepository.NewProductSubCategoryRepository(adapter.Adapters.DzikraPostgres)
	productVariantRepository := productVariantRepository.NewProductVariantRepository(adapter.Adapters.DzikraPostgres)
	productGroceryRepository := productGroceryRepository.NewProductGroceryRepository(adapter.Adapters.DzikraPostgres)
	productImageRepository := productImageRepository.NewProductImageRepository(adapter.Adapters.DzikraPostgres)

	// product  service
	productService := productService.NewProductService(
		adapter.Adapters.DzikraPostgres,
		productRepository,
		productCategoryRepository,
		productSubCategoryRepository,
		productVariantRepository,
		productGroceryRepository,
		minioService,
		productImageRepository,
	)

	// handler
	handler.service = productService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
