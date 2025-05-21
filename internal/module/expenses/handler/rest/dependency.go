package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/ports"
	expensesRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/repository"
	expensesService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/service"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type expensesHandler struct {
	service    ports.ExpensesService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewExpensesHandler() *expensesHandler {
	var handler = new(expensesHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	expensesRepository := expensesRepository.NewExpensesRepository(adapter.Adapters.DzikraPostgres)

	// expenses service
	expensesService := expensesService.NewExpensesService(
		adapter.Adapters.DzikraPostgres,
		expensesRepository,
	)

	// handler
	handler.service = expensesService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
