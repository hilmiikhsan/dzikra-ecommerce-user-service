package rest

import (
	externalOrder "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/order"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	rajaongkirService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/service"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	addressRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/repository"
	cartRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/order/ports"
	orderService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/order/service"
	productGroceryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/repository"
	voucherRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/repository"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type orderHandler struct {
	service    ports.OrderService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewOrderHandler() *orderHandler {
	var handler = new(orderHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalOrder := &externalOrder.External{}

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// integration service
	rajaongkirService := rajaongkirService.NewRajaongkirService(redisRepository)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	voucherRepository := voucherRepository.NewVoucherRepository(adapter.Adapters.DzikraPostgres)
	addressRepository := addressRepository.NewAddressRepository(adapter.Adapters.DzikraPostgres)
	cartRepository := cartRepository.NewCartRepository(adapter.Adapters.DzikraPostgres)
	productGroceryRepository := productGroceryRepository.NewProductGroceryRepository(adapter.Adapters.DzikraPostgres)

	// order  service
	orderService := orderService.NewOrderService(
		adapter.Adapters.DzikraPostgres,
		voucherRepository,
		addressRepository,
		rajaongkirService,
		cartRepository,
		productGroceryRepository,
		externalOrder,
	)

	// handler
	handler.service = orderService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
