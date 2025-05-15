package service

import (
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/ports"
	addressPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/ports"
	cartPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/ports"
	shippingPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/shipping/ports"
)

var _ shippingPorts.ShippingService = &shippingService{}

type shippingService struct {
	cartRepository    cartPorts.CartRepository
	addressRepository addressPorts.AddressRepository
	rajaongkirService rajaongkirPorts.RajaongkirService
}

func NewShippingService(
	cartRepository cartPorts.CartRepository,
	addressRepository addressPorts.AddressRepository,
	rajaongkirService rajaongkirPorts.RajaongkirService,
) *shippingService {
	return &shippingService{
		cartRepository:    cartRepository,
		addressRepository: addressRepository,
		rajaongkirService: rajaongkirService,
	}
}
