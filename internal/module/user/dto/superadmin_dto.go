package dto

// Super Admin DTO
type CreateRolePermissionRequest struct {
	Roles          string           `json:"roles" validate:"required,max=50,min=3"`
	Description    string           `json:"desc" validate:"required,max=255"`
	AppPermissions []AppPermissions `json:"app_permissions" validate:"required,dive"`
}

type AppPermissions struct {
	Action          string `json:"action" validate:"required,max=20,role_permission_action"`
	AppPermissionID string `json:"app_permission_id" validate:"required,uuid"`
	Resource        string `json:"resource" validate:"required,max=20,resource_permission_action"`
}

type CreateRolePermissionResponse struct {
	ID                 string               `json:"id"`
	Roles              string               `json:"roles"`
	Description        string               `json:"desc"`
	RoleAppPermissions []RoleAppPermissions `json:"role_app_permission"`
}

type RoleAppPermissions struct {
	ApplicationPermissionID string                `json:"applicationperm_id"`
	ApplicationPermission   ApplicationPermission `json:"application_permission"`
}

type ApplicationPermission struct {
	ApplicationID string        `json:"application_id"`
	Permissions   []Permissions `json:"permissions"`
}

type Permissions struct {
	ID       string `json:"id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}
