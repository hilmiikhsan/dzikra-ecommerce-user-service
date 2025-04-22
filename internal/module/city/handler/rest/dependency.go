package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
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

	// city service
	cityService := cityService.NewCityService(
		redisRepository,
	)

	// handler
	handler.service = cityService
	handler.validator = validator

	return handler
}
