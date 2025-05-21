package grpc

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/product_variant"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/ports"
	productVariantRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/service"
)

type productVariantGrpcAPI struct {
	ProductVariantService ports.ProductVariantService
	product_variant.UnimplementedProductVariantServiceServer
}

func NewProductVariantGrpcAPI() *productVariantGrpcAPI {
	var handler = new(productVariantGrpcAPI)

	// repository
	productVariantRepository := productVariantRepository.NewProductVariantRepository(adapter.Adapters.DzikraPostgres)

	// service
	productVariantService := service.NewProductVariantService(
		productVariantRepository,
	)

	handler.ProductVariantService = productVariantService

	return handler
}
