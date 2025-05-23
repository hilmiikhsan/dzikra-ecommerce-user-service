package dto

// Super Admin DTO
type RolePermissionRequest struct {
	Roles          string           `json:"roles" validate:"required,max=50,min=3"`
	Description    string           `json:"desc" validate:"required,max=255"`
	AppPermissions []AppPermissions `json:"app_permissions" validate:"required,dive"`
}

type AppPermissions struct {
	Action          string `json:"action" validate:"required,max=20,role_permission_action"`
	AppPermissionID string `json:"app_permission_id" validate:"required,uuid"`
	Resource        string `json:"resource" validate:"required,max=20,resource_permission_action"`
}

type RolePermissionResponse struct {
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

type GetListRole struct {
	Roles       []GetListRolePermission `json:"roles"`
	TotalPages  int                     `json:"total_page"`
	CurrentPage int                     `json:"current_page"`
	PageSize    int                     `json:"page_size"`
	TotalData   int                     `json:"total_data"`
}

type GetListRolePermission struct {
	ID                string                     `json:"id"`
	Roles             string                     `json:"roles"`
	Description       string                     `json:"desc"`
	Static            bool                       `json:"static"`
	RoleAppPermission []GetListRoleAppPermission `json:"role_app_permission"`
}

type GetListRoleAppPermission struct {
	ApplicationID string              `json:"application_id"`
	Name          string              `json:"name"`
	Permissions   []GetListPermission `json:"permission"`
}

type GetListPermission struct {
	Action                  string `json:"action"`
	ApplicationPermissionID string `json:"aplicationperm_id"`
	Resource                string `json:"resource"`
}

type GetDetailRoleResponse struct {
	ID                string                `json:"id"`
	Description       string                `json:"desc"`
	RoleAppPermission []DetailAppPermission `json:"role_app_permission"`
}

type DetailAppPermission struct {
	ApplicationID string              `json:"application_id"`
	Application   string              `json:"application"`
	Permissions   []GetListPermission `json:"permissions"`
	Roles         string              `json:"roles"`
}

type GetListApplicationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetListPermissionByAppResponse struct {
	ApplicationPermissions []PermissionApp `json:"application_permissions"`
}

type PermissionApp struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Permissions []ListPermissionApp `json:"permissions"`
}

type ListPermissionApp struct {
	AppPermiID string `json:"appperm_id"`
	Action     string `json:"action"`
	Resource   string `json:"resource"`
}

type SoftDeleteRolePermissionRequest struct {
	Roles          string           `json:"roles" validate:"required,max=50,min=3"`
	Description    string           `json:"desc" validate:"required,max=255"`
	AppPermissions []AppPermissions `json:"app_permissions" validate:"required,dive"`
}

type UserRoleDetail struct {
	ApplicationPermission []ApplicationPermissionDetail `json:"application_permission"`
	Roles                 string                        `json:"roles"`
}

type ApplicationPermissionDetail struct {
	ApplicationID string   `json:"application_id"`
	Name          string   `json:"name"`
	Permissions   []string `json:"permissions"`
}
