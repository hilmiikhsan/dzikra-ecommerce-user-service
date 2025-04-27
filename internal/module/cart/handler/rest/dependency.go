package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/ports"
	cartRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/repository"
	cartService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/service"
	productRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/repository"
	productVariantRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/repository"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type cartHandler struct {
	service    ports.CartService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewCartHandler() *cartHandler {
	var handler = new(cartHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	cartRepository := cartRepository.NewCartRepository(adapter.Adapters.DzikraPostgres)
	productRepository := productRepository.NewProductRepository(adapter.Adapters.DzikraPostgres)
	productVariantRepository := productVariantRepository.NewProductVariantRepository(adapter.Adapters.DzikraPostgres)

	// cart  service
	cartService := cartService.NewCartService(
		adapter.Adapters.DzikraPostgres,
		cartRepository,
		productRepository,
		productVariantRepository,
	)

	// handler
	handler.service = cartService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
