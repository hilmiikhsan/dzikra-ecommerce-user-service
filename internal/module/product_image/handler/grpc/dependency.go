package grpc

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/product_image"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/ports"
	productImageRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/service"
)

type productImageGrpcAPI struct {
	ProductImageService ports.ProductImageService
	product_image.UnimplementedProductImageServiceServer
}

func NewProductImageGrpcAPI() *productImageGrpcAPI {
	var handler = new(productImageGrpcAPI)

	// repository
	productImageRepository := productImageRepository.NewProductImageRepository(adapter.Adapters.DzikraPostgres)

	// service
	productImageService := service.NewProductImageService(
		adapter.Adapters.DzikraPostgres,
		productImageRepository,
	)

	handler.ProductImageService = productImageService

	return handler
}
