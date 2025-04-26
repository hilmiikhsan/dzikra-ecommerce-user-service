package service

import (
	redisPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis/ports"
	subDistrictPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/ports"
)

var _ subDistrictPorts.SubDistrictService = &subDistrictService{}

type subDistrictService struct {
	redisRepository redisPorts.RedisRepository
}

func NewSubDistrictService(
	redisRepository redisPorts.RedisRepository,
) *subDistrictService {
	return &subDistrictService{
		redisRepository: redisRepository,
	}
}
