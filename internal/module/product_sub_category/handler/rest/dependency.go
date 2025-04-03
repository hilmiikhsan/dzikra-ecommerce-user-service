package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	productCategoryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/ports"
	productSubCategoryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/repository"
	productSubCategoryService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/service"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type productSubCategoryHandler struct {
	service    ports.ProductSubCategoryService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewProductSubCategoryHandler() *productSubCategoryHandler {
	var handler = new(productSubCategoryHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	productSubCategoryRepository := productSubCategoryRepository.NewProductSubCategoryRepository(adapter.Adapters.DzikraPostgres)
	productCategoryRepository := productCategoryRepository.NewProductCategoryRepository(adapter.Adapters.DzikraPostgres)

	// product sub category service
	productSubCategoryService := productSubCategoryService.NewProductSubCategoryService(
		adapter.Adapters.DzikraPostgres,
		productSubCategoryRepository,
		productCategoryRepository,
	)

	// handler
	handler.service = productSubCategoryService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
