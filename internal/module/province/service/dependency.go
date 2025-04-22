package service

import (
	redisPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis/ports"
	provincePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/ports"
)

var _ provincePorts.ProvinceService = &provinceService{}

type provinceService struct {
	redisRepository redisPorts.RedisRepository
}

func NewProvinceService(
	redisRepository redisPorts.RedisRepository,
) *provinceService {
	return &provinceService{
		redisRepository: redisRepository,
	}
}
