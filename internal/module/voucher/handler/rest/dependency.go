package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/ports"
	voucherRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/service"
	voucherTypeRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_type/repository"
	voucherUsageRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_usage/repository"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type voucherHandler struct {
	service    ports.VoucherService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewVoucherHandler() *voucherHandler {
	var handler = new(voucherHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	voucherRepository := voucherRepository.NewVoucherRepository(adapter.Adapters.DzikraPostgres)
	voucherTypeRepository := voucherTypeRepository.NewVoucherTypeRepository(adapter.Adapters.DzikraPostgres)
	voucherUsageRepository := voucherUsageRepository.NewVoucherUsageRepository(adapter.Adapters.DzikraPostgres)

	// service
	userService := service.NewVoucherService(
		adapter.Adapters.DzikraPostgres,
		voucherRepository,
		voucherTypeRepository,
		voucherUsageRepository,
	)

	// handler
	handler.service = userService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
