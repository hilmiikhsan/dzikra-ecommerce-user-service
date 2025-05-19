package service

import (
	externalOrder "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/order"
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/ports"
	addressPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/ports"
	cartPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/ports"
	orderPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/order/ports"
	productGroceryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/ports"
	voucherPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/ports"
	"github.com/jmoiron/sqlx"
)

var _ orderPorts.OrderService = &orderService{}

type orderService struct {
	db                       *sqlx.DB
	voucherRepository        voucherPorts.VoucherRepository
	addressRepository        addressPorts.AddressRepository
	rajaongkirService        rajaongkirPorts.RajaongkirService
	cartRepository           cartPorts.CartRepository
	productGroceryRepository productGroceryPorts.ProductGroceryRepository
	externalOrder            externalOrder.ExternalOrder
}

func NewOrderService(
	db *sqlx.DB,
	voucherRepository voucherPorts.VoucherRepository,
	addressRepository addressPorts.AddressRepository,
	rajaongkirService rajaongkirPorts.RajaongkirService,
	cartRepository cartPorts.CartRepository,
	productGroceryRepository productGroceryPorts.ProductGroceryRepository,
	externalOrder externalOrder.ExternalOrder,
) *orderService {
	return &orderService{
		db:                       db,
		voucherRepository:        voucherRepository,
		addressRepository:        addressRepository,
		rajaongkirService:        rajaongkirService,
		cartRepository:           cartRepository,
		productGroceryRepository: productGroceryRepository,
		externalOrder:            externalOrder,
	}
}
