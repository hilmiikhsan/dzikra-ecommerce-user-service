package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/service"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/ports"
	provinceService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/service"
)

type provinceHandler struct {
	service   ports.ProvinceService
	validator adapter.Validator
}

func NewProvinceHandler() *provinceHandler {
	var handler = new(provinceHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// integration service
	rajaongkirService := rajaongkirPorts.NewRajaongkirService(redisRepository)

	// province service
	provinceService := provinceService.NewProvinceService(
		rajaongkirService,
	)

	// handler
	handler.service = provinceService
	handler.validator = validator

	return handler
}
