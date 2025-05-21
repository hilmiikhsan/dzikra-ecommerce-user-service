package rest

import (
	externalOrder "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/order"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/dashboard/ports"
	dashboardService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/dashboard/service"
	expensesRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/repository"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type dashboardHandler struct {
	service    ports.DashboardService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewDashboardHandler() *dashboardHandler {
	var handler = new(dashboardHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalOrder := &externalOrder.External{}

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	expensesRepository := expensesRepository.NewExpensesRepository(adapter.Adapters.DzikraPostgres)

	// dashboard service
	dashboardService := dashboardService.NewDashboardService(
		adapter.Adapters.DzikraPostgres,
		expensesRepository,
		externalOrder,
	)

	// handler
	handler.service = dashboardService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
