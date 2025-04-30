package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/service"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/ports"
	cityService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/service"
)

type cityHandler struct {
	service   ports.CityService
	validator adapter.Validator
}

func NewCityHandler() *cityHandler {
	var handler = new(cityHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// integration service
	rajaongkirService := rajaongkirPorts.NewRajaongkirService(redisRepository)

	// city service
	cityService := cityService.NewCityService(
		rajaongkirService,
	)

	// handler
	handler.service = cityService
	handler.validator = validator

	return handler
}
