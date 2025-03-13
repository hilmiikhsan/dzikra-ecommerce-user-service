package rest

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-user-service/external/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/middleware"
	applicationPermissionRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/application_permission/repository"
	applicationRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/list_application/repository"
	roleRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/repository"
	roleAppPermissionRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_app_permission/repository"
	rolePermissionRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_permission/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/ports"
	userRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/service"
	userFcmTokenRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_fcm_token/repository"
	userProfileRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/repository"
	userRoleRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/repository"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/jwt_handler"
)

type userHandler struct {
	service    ports.UserService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

type superAdminHandler struct {
	service    ports.SuperAdminService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewUserHandler() *userHandler {
	var handler = new(userHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalNotification := &externalNotification.External{}

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	userRepository := userRepository.NewUserRepository(adapter.Adapters.DzikraPostgres)
	roleRepository := roleRepository.NewRoleRepository(adapter.Adapters.DzikraPostgres)
	userRoleRepository := userRoleRepository.NewUserRoleRepository(adapter.Adapters.DzikraPostgres)
	userProfileRepository := userProfileRepository.NewUserProfileRepository(adapter.Adapters.DzikraPostgres)
	rolePermissionRepository := rolePermissionRepository.NewRolePermissionRepository(adapter.Adapters.DzikraPostgres)
	userFcmTokenRepository := userFcmTokenRepository.NewUserFcmTokenRepository(adapter.Adapters.DzikraPostgres)

	// service
	userService := service.NewUserService(
		adapter.Adapters.DzikraPostgres,
		userRepository,
		roleRepository,
		userRoleRepository,
		userProfileRepository,
		redisRepository,
		externalNotification,
		jwt,
		rolePermissionRepository,
		userFcmTokenRepository,
	)

	// handler
	handler.service = userService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}

func NewSuperAdminHandler() *superAdminHandler {
	var handler = new(superAdminHandler)

	// validator
	validator := adapter.Adapters.Validator

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// repository
	applicationPermissionRepository := applicationPermissionRepository.NewApplicationPermissionRepository(adapter.Adapters.DzikraPostgres)
	roleRepository := roleRepository.NewRoleRepository(adapter.Adapters.DzikraPostgres)
	roleAppPermissionRepository := roleAppPermissionRepository.NewRoleAppPermissionRepository(adapter.Adapters.DzikraPostgres)
	applicationRepository := applicationRepository.NewApplicationRepository(adapter.Adapters.DzikraPostgres)
	rolePermissionRepository := rolePermissionRepository.NewRolePermissionRepository(adapter.Adapters.DzikraPostgres)
	userRoleRepository := userRoleRepository.NewUserRoleRepository(adapter.Adapters.DzikraPostgres)

	// service
	superAdminService := service.NewSuperAdminService(
		adapter.Adapters.DzikraPostgres,
		applicationPermissionRepository,
		roleRepository,
		roleAppPermissionRepository,
		applicationRepository,
		rolePermissionRepository,
		userRoleRepository,
	)

	// handler
	handler.service = superAdminService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
