package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/list_application/ports"
	applicationRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/list_application/repository"
	applicationService "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/list_application/service"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/jwt_handler"
)

type applicationHandler struct {
	service    ports.ApplicationService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewApplicationHandler() *applicationHandler {
	var handler = new(applicationHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	applicationRepository := applicationRepository.NewApplicationRepository(adapter.Adapters.DzikraPostgres)

	// role service
	applicationService := applicationService.NewApplicationService(
		adapter.Adapters.DzikraPostgres,
		applicationRepository,
	)

	// handler
	handler.service = applicationService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
