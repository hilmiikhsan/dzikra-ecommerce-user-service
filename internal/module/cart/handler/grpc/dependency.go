package grpc

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/cart"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/ports"
	cartRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/service"
	productRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/repository"
	productVariantRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/repository"
)

type cartGrpcAPI struct {
	CartService ports.CartService
	cart.UnimplementedCartServiceServer
}

func NewCartGrpcAPI() *cartGrpcAPI {
	var handler = new(cartGrpcAPI)

	// repository
	cartRepository := cartRepository.NewCartRepository(adapter.Adapters.DzikraPostgres)
	productRepository := productRepository.NewProductRepository(adapter.Adapters.DzikraPostgres)
	productVariantRepository := productVariantRepository.NewProductVariantRepository(adapter.Adapters.DzikraPostgres)

	// service
	cartService := service.NewCartService(
		adapter.Adapters.DzikraPostgres,
		cartRepository,
		productRepository,
		productVariantRepository,
	)

	handler.CartService = cartService

	return handler
}
