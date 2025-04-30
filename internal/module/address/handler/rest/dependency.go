package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	rajaongkirService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/service"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/ports"
	addressRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/repository"
	addressService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/service"
	cityService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/service"
	provinceService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/service"
	subDistrictPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/service"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type addressHandler struct {
	service    ports.AddressService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewAddressHandler() *addressHandler {
	var handler = new(addressHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	addressRepository := addressRepository.NewAddressRepository(adapter.Adapters.DzikraPostgres)

	// integration service
	rajaongkirService := rajaongkirService.NewRajaongkirService(redisRepository)

	// service
	provinceService := provinceService.NewProvinceService(rajaongkirService)
	cityService := cityService.NewCityService(rajaongkirService)
	subDistrictPorts := subDistrictPorts.NewSubDistrictService(rajaongkirService)

	// address service
	addressService := addressService.NewAddressService(
		adapter.Adapters.DzikraPostgres,
		addressRepository,
		redisRepository,
		provinceService,
		cityService,
		subDistrictPorts,
	)

	// handler
	handler.service = addressService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
