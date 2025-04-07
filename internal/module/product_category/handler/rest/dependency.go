package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/ports"
	productCategoryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/repository"
	productCategoryService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/service"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type productCategoryHandler struct {
	service    ports.ProductCategoryService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewProductCategoryHandler() *productCategoryHandler {
	var handler = new(productCategoryHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	productCategoryRepository := productCategoryRepository.NewProductCategoryRepository(adapter.Adapters.DzikraPostgres)

	// product category service
	productCategoryService := productCategoryService.NewProductCategoryService(
		productCategoryRepository,
	)

	// handler
	handler.service = productCategoryService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
