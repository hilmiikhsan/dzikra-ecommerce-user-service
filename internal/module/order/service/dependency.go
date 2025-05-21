package service

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/notification"
	externalOrder "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/order"
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/ports"
	addressPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/ports"
	cartPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/ports"
	orderPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/order/ports"
	productPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/ports"
	productGroceryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/ports"
	productVariantPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/ports"
	userFcmTokenPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user_fcm_token/ports"
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
	productVariantRepository productVariantPorts.ProductVariantRepository
	productRepository        productPorts.ProductRepository
	userFcmTokenRepository   userFcmTokenPorts.UserFCMTokenRepository
	externalNotification     externalNotification.ExternalNotification
}

func NewOrderService(
	db *sqlx.DB,
	voucherRepository voucherPorts.VoucherRepository,
	addressRepository addressPorts.AddressRepository,
	rajaongkirService rajaongkirPorts.RajaongkirService,
	cartRepository cartPorts.CartRepository,
	productGroceryRepository productGroceryPorts.ProductGroceryRepository,
	externalOrder externalOrder.ExternalOrder,
	productVariantRepository productVariantPorts.ProductVariantRepository,
	productRepository productPorts.ProductRepository,
	userFcmTokenRepository userFcmTokenPorts.UserFCMTokenRepository,
	externalNotification externalNotification.ExternalNotification,
) *orderService {
	return &orderService{
		db:                       db,
		voucherRepository:        voucherRepository,
		addressRepository:        addressRepository,
		rajaongkirService:        rajaongkirService,
		cartRepository:           cartRepository,
		productGroceryRepository: productGroceryRepository,
		externalOrder:            externalOrder,
		productVariantRepository: productVariantRepository,
		productRepository:        productRepository,
		userFcmTokenRepository:   userFcmTokenRepository,
		externalNotification:     externalNotification,
	}
}
