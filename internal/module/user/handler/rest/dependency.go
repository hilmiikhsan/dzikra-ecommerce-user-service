package rest

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-user-service/external/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/redis"
	roleRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/repository"
	userRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/service"
	userProfileRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/repository"
	userRoleRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/repository"
)

func NewUserHandler() *userHandler {
	var handler = new(userHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalNotification := &externalNotification.External{}

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	// jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	// middlewareHandler := middleware.NewAuthMiddleware(jwt)

	// repository
	userRepository := userRepository.NewUserRepository(adapter.Adapters.DzikraPostgres)
	roleRepository := roleRepository.NewRoleRepository(adapter.Adapters.DzikraPostgres)
	userRoleRepository := userRoleRepository.NewUserRoleRepository(adapter.Adapters.DzikraPostgres)
	userProfileRepository := userProfileRepository.NewUserProfileRepository(adapter.Adapters.DzikraPostgres)

	// service
	userService := service.NewUserService(
		adapter.Adapters.DzikraPostgres,
		userRepository,
		roleRepository,
		userRoleRepository,
		userProfileRepository,
		redisRepository,
		externalNotification,
	)

	// handler
	handler.service = userService
	// handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
