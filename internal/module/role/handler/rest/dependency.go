package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	applicationPermissionRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/application_permission/repository"
	applicationRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role/ports"
	roleRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role/repository"
	roleService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role/service"
	roleAppPermissionRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role_app_permission/repository"
	rolePermissionRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role_permission/repository"
	userRoleRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user_role/repository"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type roleHandler struct {
	service    ports.RoleService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewRoleHandler() *roleHandler {
	var handler = new(roleHandler)

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

	// role service
	roleService := roleService.NewRoleService(
		adapter.Adapters.DzikraPostgres,
		applicationPermissionRepository,
		roleRepository,
		roleAppPermissionRepository,
		applicationRepository,
		rolePermissionRepository,
		userRoleRepository,
	)

	// handler
	handler.service = roleService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
