package service

import (
	applicationPermissionPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/application_permission/ports"
	applicationPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/ports"
	rolePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role/ports"
	roleAppPermissionPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role_app_permission/ports"
	rolePermissionPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role_permission/ports"
	userRolePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user_role/ports"
	"github.com/jmoiron/sqlx"
)

var _ rolePorts.RoleService = &roleService{}

type roleService struct {
	db                              *sqlx.DB
	applicationPermissionRepository applicationPermissionPorts.ApplicationPermissionRepository
	roleRepository                  rolePorts.RoleRepository
	roleAppPermissionRepository     roleAppPermissionPorts.RoleAppPermissionRepository
	applicationRepository           applicationPorts.ApplicationRepository
	rolePermissionRepository        rolePermissionPorts.RolePermissionRepository
	userRoleRepository              userRolePorts.UserRoleRepository
}

func NewRoleService(
	db *sqlx.DB,
	applicationPermissionRepository applicationPermissionPorts.ApplicationPermissionRepository,
	roleRepository rolePorts.RoleRepository,
	roleAppPermissionRepository roleAppPermissionPorts.RoleAppPermissionRepository,
	applicationRepository applicationPorts.ApplicationRepository,
	rolePermissionRepository rolePermissionPorts.RolePermissionRepository,
	userRoleRepository userRolePorts.UserRoleRepository,
) *roleService {
	return &roleService{
		db:                              db,
		applicationPermissionRepository: applicationPermissionRepository,
		roleRepository:                  roleRepository,
		roleAppPermissionRepository:     roleAppPermissionRepository,
		applicationRepository:           applicationRepository,
		rolePermissionRepository:        rolePermissionRepository,
		userRoleRepository:              userRoleRepository,
	}
}
