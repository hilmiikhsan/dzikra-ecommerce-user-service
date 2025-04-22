package service

import (
	redisPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis/ports"
	cityPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/ports"
)

var _ cityPorts.CityService = &cityService{}

type cityService struct {
	redisRepository redisPorts.RedisRepository
}

func NewCityService(
	redisRepository redisPorts.RedisRepository,
) *cityService {
	return &cityService{
		redisRepository: redisRepository,
	}
}
