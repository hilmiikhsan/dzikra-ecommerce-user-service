package service

import (
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/ports"
	subDistrictPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/ports"
)

var _ subDistrictPorts.SubDistrictService = &subDistrictService{}

type subDistrictService struct {
	rajaongkirService rajaongkirPorts.RajaongkirService
}

func NewSubDistrictService(
	rajaongkirService rajaongkirPorts.RajaongkirService,
) *subDistrictService {
	return &subDistrictService{
		rajaongkirService: rajaongkirService,
	}
}
