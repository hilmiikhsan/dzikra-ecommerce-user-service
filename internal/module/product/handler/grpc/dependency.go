package grpc

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/product"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/ports"
	productRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/service"
)

type productGrpcAPI struct {
	ProductService ports.ProductService
	product.UnimplementedProductServiceServer
}

func NewProductGrpcAPI() *productGrpcAPI {
	var handler = new(productGrpcAPI)

	// repository
	productRepository := productRepository.NewProductRepository(adapter.Adapters.DzikraPostgres)

	// service
	productService := service.NewProductService(
		nil,
		productRepository,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	handler.ProductService = productService

	return handler
}
