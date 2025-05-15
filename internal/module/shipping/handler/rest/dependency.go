package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/service"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	addressRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/repository"
	cartRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/shipping/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/shipping/service"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type shippingHandler struct {
	service    ports.ShippingService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewShippingHandler() *shippingHandler {
	var handler = new(shippingHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// integration service
	rajaongkirService := rajaongkirPorts.NewRajaongkirService(redisRepository)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	cartRepository := cartRepository.NewCartRepository(adapter.Adapters.DzikraPostgres)
	addressRepository := addressRepository.NewAddressRepository(adapter.Adapters.DzikraPostgres)

	// shipping service
	shippingService := service.NewShippingService(
		cartRepository,
		addressRepository,
		rajaongkirService,
	)

	// handler
	handler.service = shippingService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
