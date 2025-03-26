package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/entity"
	roleApppermission "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_app_permission/entity"
	rolePermission "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_permission/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *roleService) GetDetailRole(ctx context.Context, roleID string) (*dto.GetDetailRoleResponse, error) {
	// find role by ID
	roleResult, err := s.roleRepository.FindRoleByID(ctx, roleID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrRoleNotFound) {
			log.Error().Any("roleID", roleID).Msg("service::GetDetailRole - Role not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrRoleNotFound))
		}
		log.Error().Err(err).Any("roleID", roleID).Msg("service::GetDetailRole - Failed to get role by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check role
	if roleResult == nil {
		log.Error().Any("roleID", roleID).Msg("service::GetDetailRole - Role not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrRoleNotFound))
	}

	// mapping data role app permission
	converted := make([]dto.DetailAppPermission, 0)
	if roleResult.RoleAppPermission != nil {
		for _, data := range roleResult.RoleAppPermission {
			converted = append(converted, utils.MapToDetailAppPermission(data, roleResult.Roles))
		}
	}

	// mapping get detail role response data
	response := &dto.GetDetailRoleResponse{
		ID:                roleResult.ID,
		Description:       roleResult.Description,
		RoleAppPermission: converted,
	}

	// return response
	return response, nil
}

func (s *roleService) GetListRole(ctx context.Context, page, limit int, search string) (*dto.GetListRole, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list role
	roleResults, total, err := s.roleRepository.FindListRole(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListRole - Failed to get list role")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if roleResults is nil
	if roleResults == nil {
		roleResults = []dto.GetListRolePermission{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create mapping response
	responses := &dto.GetListRole{
		Roles:       roleResults,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return responses, nil
}

func (s *roleService) CreateRolePermission(ctx context.Context, req *dto.RolePermissionRequest) (*dto.RolePermissionResponse, error) {
	// mapping request apllication permission ID
	applicationIDs := make([]string, 0, len(req.AppPermissions))
	for _, item := range req.AppPermissions {
		applicationIDs = append(applicationIDs, item.AppPermissionID)
	}

	// get application permissions by ID
	isExist, err := s.applicationPermissionRepository.FindApplicationPermissionByID(ctx, applicationIDs)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateRolePermission - Failed to get application permissions")
		return nil, err
	}

	// check if application permissions not found
	if !isExist {
		log.Error().Msg("service::CreateRolePermission - Application permissions not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrApplicationPermissionNotFound))
	}

	// get role data by name
	roleResult, err := s.roleRepository.FindRoleByName(ctx, strings.ToUpper(req.Roles))
	if err != nil {
		if err.Error() == constants.ErrRoleNotFound {
			log.Error().Any("roleName", req.Roles).Msg("service::CreateRolePermission - Role not found")
		} else {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateRolePermission - Failed to find role")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	// check if role already exist
	if roleResult != nil {
		log.Error().Any("roleName", req.Roles).Msg("service::CreateRolePermission - Role already exist")
		return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrRoleAlreadyExist))
	}

	// begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateRolePermission - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateRolePermission - Failed to rollback transaction")
			}
		}
	}()

	// generate role UUID
	generateRoleID, err := utils.GenerateUUIDv7String()
	if err != nil {
		log.Error().Err(err).Msg("service::CreateRolePermission - Failed to generate role UUID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// insert new role
	err = s.roleRepository.InsertNewRole(ctx, tx, &entity.Role{
		ID:          generateRoleID,
		Name:        strings.ToUpper(req.Roles),
		Description: req.Description,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateRolePermission - Failed to insert new role")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping request role app permissions
	var roleAppPermissions []roleApppermission.RoleAppPermission
	for _, item := range req.AppPermissions {
		roleAppPermissionID, err := utils.GenerateUUIDv7String()
		if err != nil {
			log.Error().Err(err).Msg("service::CreateRolePermission - Failed to generate role app permission UUID")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		appPermUUID, _ := uuid.Parse(item.AppPermissionID)

		roleAppPermissions = append(roleAppPermissions, roleApppermission.RoleAppPermission{
			ID:              roleAppPermissionID,
			RoleID:          generateRoleID,
			AppPermissionID: appPermUUID,
		})
	}

	// Insert role app permissions (bulk insert)
	err = s.roleAppPermissionRepository.InsertNewRoleAppPermissions(ctx, tx, roleAppPermissions)
	if err != nil {
		log.Error().Err(err).Any("payload", roleAppPermissions).Msg("service::CreateRolePermission - Failed to insert new role app permissions")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping role permissions
	var rolePermissions []rolePermission.RolePermission
	for _, item := range req.AppPermissions {
		permissionID, err := s.applicationPermissionRepository.GetPermissionIDByID(ctx, item.AppPermissionID)
		if err != nil {
			log.Error().Err(err).Msg("service::CreateRolePermission - Failed to get permission ID for app permission")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		newRolePermissionID, err := utils.GenerateUUIDv7String()
		if err != nil {
			log.Error().Err(err).Msg("service::CreateRolePermission - Failed to generate role permission UUID")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		permissionUUID, _ := uuid.Parse(permissionID)

		rolePermissions = append(rolePermissions, rolePermission.RolePermission{
			ID:           newRolePermissionID,
			RoleID:       generateRoleID,
			PermissionID: permissionUUID,
		})
	}

	// Insert role permissions (bulk insert)
	err = s.rolePermissionRepository.InsertNewRolePermissions(ctx, tx, rolePermissions)
	if err != nil {
		log.Error().Err(err).Any("payload", rolePermissions).Msg("service::CreateRolePermission - Failed to insert new role permissions")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateRolePermission - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get role permissions
	response, err := s.roleRepository.FindRolePermission(ctx, generateRoleID.String())
	if err != nil {
		log.Error().Err(err).Any("roleID", generateRoleID).Msg("service::CreateRolePermission - Failed to get role permissions")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// return response
	return response, nil
}

func (s *roleService) RemoveRolePermission(ctx context.Context, roleID string) error {
	// begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::RemoveRolePermission - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Str("roleID", roleID).Msg("service::RemoveRolePermission - Failed to rollback transaction")
			}
		}
	}()

	// check if role exist
	roleResult, err := s.roleRepository.FindRoleByID(ctx, roleID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrRoleNotFound) {
			log.Error().Str("roleID", roleID).Msg("service::RemoveRolePermission - Role not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrRoleNotFound))
		}

		log.Error().Err(err).Str("roleID", roleID).Msg("service::RemoveRolePermission - Failed to find role")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check value static role
	if roleResult.Static {
		log.Error().Str("roleID", roleID).Msg("service::RemoveRolePermission - Static role cannot be deleted")
		return err_msg.NewCustomErrors(fiber.StatusForbidden, err_msg.WithMessage(constants.ErrStaticRoleCannotBeDeleted))
	}

	// soft delete role app permissions
	err = s.roleAppPermissionRepository.SoftDeleteRoleAppPermissions(ctx, tx, roleResult.ID)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::RemoveRolePermission - Failed to remove role app permissions")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// soft delete role permissions
	err = s.rolePermissionRepository.SoftDeleteRolePermissions(ctx, tx, roleResult.ID)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::RemoveRolePermission - Failed to remove role permissions")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// soft delete user roles
	err = s.userRoleRepository.SoftDeleteUserRolePermissions(ctx, tx, roleResult.ID)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::RemoveRolePermission - Failed to remove user roles")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// sopft delete roles
	err = s.roleRepository.SoftDeleteRole(ctx, tx, roleResult.ID)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::RemoveRolePermission - Failed to remove role")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::RemoveRolePermission - Failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// return nil
	return nil
}

func (s *roleService) UpdateRolePermission(ctx context.Context, req *dto.SoftDeleteRolePermissionRequest, roleID string) (*dto.RolePermissionResponse, error) {
	// mapping request apllication permission ID
	applicationIDs := make([]string, 0, len(req.AppPermissions))
	for _, item := range req.AppPermissions {
		applicationIDs = append(applicationIDs, item.AppPermissionID)
	}

	// get application permissions by ID
	isExist, err := s.applicationPermissionRepository.FindApplicationPermissionByID(ctx, applicationIDs)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateUserRoles - Failed to get application permissions")
		return nil, err
	}

	// check if application permissions not found
	if !isExist {
		log.Error().Msg("service::UpdateUserRoles - Application permissions not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrApplicationPermissionNotFound))
	}

	// get role data by ID
	existingRole, err := s.roleRepository.FindRoleByID(ctx, roleID)
	if err != nil {
		if err.Error() == constants.ErrRoleNotFound {
			log.Error().Str("roleID", roleID).Msg("service::UpdateRolePermission - Role not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrRoleNotFound))
		}

		log.Error().Err(err).Str("roleID", roleID).Msg("service::UpdateRolePermission - Failed to get role by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if role not found
	if existingRole == nil {
		log.Error().Str("roleID", roleID).Msg("service::UpdateRolePermission - Role not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrRoleNotFound))
	}

	// get role data by name
	newRoleName := strings.ToUpper(req.Roles)
	if newRoleName != existingRole.Roles {
		conflictRole, err := s.roleRepository.FindRoleByName(ctx, newRoleName)
		if err != nil && err.Error() != constants.ErrRoleNotFound {
			log.Error().Err(err).Str("newRoleName", newRoleName).Msg("service::UpdateRolePermission - Error checking role conflict")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		if conflictRole != nil {
			log.Error().Str("newRoleName", newRoleName).Msg("service::UpdateRolePermission - Role already exist")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrRoleAlreadyExist))
		}
	}

	// begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateUserRoles - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::UpdateUserRoles - Failed to rollback transaction")
			}
		}
	}()

	// update new role data
	err = s.roleRepository.UpdateRole(ctx, tx, roleID, newRoleName, req.Description, existingRole.Roles)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::UpdateRolePermission - Failed to update role")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Soft delete data role_app_permissions
	err = s.roleAppPermissionRepository.SoftDeleteRoleAppPermissions(ctx, tx, roleID)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::UpdateRolePermission - Failed to soft delete role app permissions")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Soft delete data role_permissions
	err = s.rolePermissionRepository.SoftDeleteRolePermissions(ctx, tx, roleID)
	if err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::UpdateRolePermission - Failed to soft delete role permissions")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Mapping new data role app permissions
	var newRoleAppPermissions []roleApppermission.RoleAppPermission
	for _, item := range req.AppPermissions {
		roleAppPermissionID, err := utils.GenerateUUIDv7String()
		if err != nil {
			log.Error().Err(err).Msg("service::UpdateRolePermission - Failed to generate role app permission UUID")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		appPermUUID, _ := uuid.Parse(item.AppPermissionID)
		roleUUID, _ := uuid.Parse(roleID)

		newRoleAppPermissions = append(newRoleAppPermissions, roleApppermission.RoleAppPermission{
			ID:              roleAppPermissionID,
			RoleID:          roleUUID,
			AppPermissionID: appPermUUID,
		})
	}

	// insert new role app permissions
	err = s.roleAppPermissionRepository.InsertNewRoleAppPermissions(ctx, tx, newRoleAppPermissions)
	if err != nil {
		log.Error().Err(err).Any("payload", newRoleAppPermissions).Msg("service::UpdateRolePermission - Failed to insert new role app permissions")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Mapping new data role permissions
	var newRolePermissions []rolePermission.RolePermission
	for _, item := range req.AppPermissions {
		permissionID, err := s.applicationPermissionRepository.GetPermissionIDByID(ctx, item.AppPermissionID)
		if err != nil {
			log.Error().Err(err).Msg("service::UpdateRolePermission - Failed to get permission ID for app permission")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		newRolePermissionID, err := utils.GenerateUUIDv7String()
		if err != nil {
			log.Error().Err(err).Msg("service::UpdateRolePermission - Failed to generate role permission UUID")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		roleUUID, _ := uuid.Parse(roleID)
		permissionUUID, _ := uuid.Parse(permissionID)

		newRolePermissions = append(newRolePermissions, rolePermission.RolePermission{
			ID:           newRolePermissionID,
			RoleID:       roleUUID,
			PermissionID: permissionUUID,
		})
	}

	// insert new data role permissions
	err = s.rolePermissionRepository.InsertNewRolePermissions(ctx, tx, newRolePermissions)
	if err != nil {
		log.Error().Err(err).Any("payload", newRolePermissions).Msg("service::UpdateRolePermission - Failed to insert new role permissions")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Str("roleID", roleID).Msg("service::UpdateRolePermission - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get new role permissions
	response, err := s.roleRepository.FindRolePermission(ctx, roleID)
	if err != nil {
		log.Error().Err(err).Any("roleID", roleID).Msg("service::UpdateRolePermission - Failed to get updated role permissions")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return response, nil
}
