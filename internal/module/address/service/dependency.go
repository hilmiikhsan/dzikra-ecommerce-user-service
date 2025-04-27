package service

import (
	redisPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis/ports"
	addressPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/ports"
	cityPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/ports"
	provincePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/ports"
	subDistrictPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/ports"
	"github.com/jmoiron/sqlx"
)

var _ addressPorts.AddressService = &addressService{}

type addressService struct {
	db                 *sqlx.DB
	addressRepository  addressPorts.AddressRepository
	redisRepository    redisPorts.RedisRepository
	provinceService    provincePorts.ProvinceService
	cityService        cityPorts.CityService
	subDistrictService subDistrictPorts.SubDistrictService
}

func NewAddressService(
	db *sqlx.DB,
	addressRepository addressPorts.AddressRepository,
	redisRepository redisPorts.RedisRepository,
	provinceService provincePorts.ProvinceService,
	cityService cityPorts.CityService,
	subDistrictService subDistrictPorts.SubDistrictService,
) *addressService {
	return &addressService{
		db:                 db,
		addressRepository:  addressRepository,
		redisRepository:    redisRepository,
		provinceService:    provinceService,
		cityService:        cityService,
		subDistrictService: subDistrictService,
	}
}
