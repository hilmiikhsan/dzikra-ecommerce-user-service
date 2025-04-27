package service

import (
	cartPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/ports"
	productPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/ports"
	productVariantPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/ports"
	"github.com/jmoiron/sqlx"
)

var _ cartPorts.CartService = &cartService{}

type cartService struct {
	db                       *sqlx.DB
	cartRepository           cartPorts.CartRepository
	productRepository        productPorts.ProductRepository
	productVariantRepository productVariantPorts.ProductVariantRepository
}

func NewCartService(
	db *sqlx.DB,
	cartRepository cartPorts.CartRepository,
	productRepository productPorts.ProductRepository,
	productVariantRepository productVariantPorts.ProductVariantRepository,
) *cartService {
	return &cartService{
		db:                       db,
		cartRepository:           cartRepository,
		productRepository:        productRepository,
		productVariantRepository: productVariantRepository,
	}
}
