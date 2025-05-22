package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/storage/minio"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/ports"
	bannerRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/repository"
	bannerService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/service"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type bannerHandler struct {
	service    ports.BannerService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewBannerHandler() *bannerHandler {
	var handler = new(bannerHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// minio service
	minioService := minio.NewMinioService(adapter.Adapters.DzikraMinio, config.Envs.MinioStorage.Bucket)

	// repository
	bannerRepository := bannerRepository.NewBannerRepository(adapter.Adapters.DzikraPostgres)

	// product  service
	bannerService := bannerService.NewBannerService(
		bannerRepository,
		minioService,
	)

	// handler
	handler.service = bannerService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
