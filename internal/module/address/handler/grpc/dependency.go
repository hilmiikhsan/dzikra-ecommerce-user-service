package grpc

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/address"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/ports"
	addressRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/service"
)

type addressImageGrpcAPI struct {
	AddressService ports.AddressService
	address.UnimplementedAddressServiceServer
}

func NewAddressGrpcAPI() *addressImageGrpcAPI {
	var handler = new(addressImageGrpcAPI)

	// repository
	addressRepository := addressRepository.NewAddressRepository(adapter.Adapters.DzikraPostgres)

	// service
	addressService := service.NewAddressService(
		adapter.Adapters.DzikraPostgres,
		addressRepository,
		nil,
		nil,
		nil,
		nil,
	)

	handler.AddressService = addressService

	return handler
}
