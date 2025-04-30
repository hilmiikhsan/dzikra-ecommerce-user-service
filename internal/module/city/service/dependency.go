package service

import (
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/ports"
	cityPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/ports"
)

var _ cityPorts.CityService = &cityService{}

type cityService struct {
	rajaongkirService rajaongkirPorts.RajaongkirService
}

func NewCityService(
	rajaongkirService rajaongkirPorts.RajaongkirService,
) *cityService {
	return &cityService{
		rajaongkirService: rajaongkirService,
	}
}
