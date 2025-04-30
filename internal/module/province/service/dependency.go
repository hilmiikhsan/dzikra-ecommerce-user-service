package service

import (
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/ports"
	provincePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/ports"
)

var _ provincePorts.ProvinceService = &provinceService{}

type provinceService struct {
	rajaongkirService rajaongkirPorts.RajaongkirService
}

func NewProvinceService(
	rajaongkirService rajaongkirPorts.RajaongkirService,
) *provinceService {
	return &provinceService{
		rajaongkirService: rajaongkirService,
	}
}
