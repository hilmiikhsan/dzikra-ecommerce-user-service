package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/ports"
	subDistrictService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/service"
)

type subDistrict struct {
	service   ports.SubDistrictService
	validator adapter.Validator
}

func NewSubDistrict() *subDistrict {
	var handler = new(subDistrict)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// sub district service
	subDistrictService := subDistrictService.NewSubDistrictService(
		redisRepository,
	)

	// handler
	handler.service = subDistrictService
	handler.validator = validator

	return handler
}
