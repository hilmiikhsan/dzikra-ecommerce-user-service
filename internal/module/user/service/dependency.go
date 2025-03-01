package service

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-user-service/external/notification"
	redisPorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/redis/ports"
	rolePorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/ports"
	rolePermissionPorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_permission/ports"
	userPorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/ports"
	userFCMTokenPorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_fcm_token/ports"
	userProfilePorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/ports"
	userRolePorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/jwt_handler"
	"github.com/jmoiron/sqlx"
)

var _ userPorts.UserService = &userService{}

type userService struct {
	db                       *sqlx.DB
	userRepository           userPorts.UserRepository
	roleRepository           rolePorts.RoleRepository
	userRoleRepository       userRolePorts.UserRoleRepository
	userProfileRepository    userProfilePorts.UserProfileRepository
	redisRepository          redisPorts.RedisRepository
	externalNotification     externalNotification.ExternalNotification
	jwt                      jwt_handler.JWT
	rolePermissionRepository rolePermissionPorts.RolePermissionRepository
	userFcmTokenRepository   userFCMTokenPorts.UserFCMTokenRepository
}

func NewUserService(db *sqlx.DB, userRepository userPorts.UserRepository, roleRepository rolePorts.RoleRepository, userRoleRepository userRolePorts.UserRoleRepository, userProfileRepository userProfilePorts.UserProfileRepository, redisRepository redisPorts.RedisRepository, externalNotification externalNotification.ExternalNotification, jwt jwt_handler.JWT, rolePermissionRepository rolePermissionPorts.RolePermissionRepository, userFcmTokenRepository userFCMTokenPorts.UserFCMTokenRepository) *userService {
	return &userService{
		db:                       db,
		userRepository:           userRepository,
		roleRepository:           roleRepository,
		userRoleRepository:       userRoleRepository,
		userProfileRepository:    userProfileRepository,
		redisRepository:          redisRepository,
		externalNotification:     externalNotification,
		jwt:                      jwt,
		rolePermissionRepository: rolePermissionRepository,
		userFcmTokenRepository:   userFcmTokenRepository,
	}
}
