package utils

import (
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role_permission/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
)

func MapUserRoleResponse(rows []entity.UserRolePermission) []dto.UserRoleDetail {
	roleMap := make(map[string]map[string]*dto.ApplicationPermissionDetail)

	for _, row := range rows {
		if !row.ApplicationID.Valid {
			continue
		}

		appID := row.ApplicationID.String
		appName := ""
		if row.ApplicationName.Valid {
			appName = row.ApplicationName.String
		}

		if _, exists := roleMap[row.RoleName]; !exists {
			roleMap[row.RoleName] = make(map[string]*dto.ApplicationPermissionDetail)
		}
		appMap := roleMap[row.RoleName]

		if _, exists := appMap[appID]; !exists {
			appMap[appID] = &dto.ApplicationPermissionDetail{
				ApplicationID: appID,
				Name:          appName,
				Permissions:   []string{},
			}
		}

		existing := appMap[appID].Permissions
		alreadyExists := false
		for _, perm := range existing {
			if perm == row.Permission {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			appMap[appID].Permissions = append(appMap[appID].Permissions, row.Permission)
		}
	}

	var results []dto.UserRoleDetail
	for roleName, appMap := range roleMap {
		var appPermissions []dto.ApplicationPermissionDetail
		for _, appPerm := range appMap {
			if appPerm.Permissions == nil {
				appPerm.Permissions = []string{}
			}
			appPermissions = append(appPermissions, *appPerm)
		}
		results = append(results, dto.UserRoleDetail{
			ApplicationPermission: appPermissions,
			Roles:                 roleName,
		})
	}

	return results
}

func MapRoleAppPermission(data dto.GetListRoleAppPermission) dto.RoleAppPermissions {
	var perms []dto.Permissions
	for _, gp := range data.Permissions {
		perms = append(perms, dto.Permissions{
			ID:       gp.ApplicationPermissionID,
			Resource: gp.Resource,
			Action:   gp.Action,
		})
	}

	var appPermID string
	if len(data.Permissions) > 0 {
		appPermID = data.Permissions[0].ApplicationPermissionID
	}

	return dto.RoleAppPermissions{
		ApplicationPermissionID: appPermID,
		ApplicationPermission: dto.ApplicationPermission{
			ApplicationID: data.ApplicationID,
			Permissions:   perms,
		},
	}
}
