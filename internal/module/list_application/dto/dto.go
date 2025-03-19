package dto

type GetListApplicationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

type GetListPermissionByAppResponse struct {
	ApplicationPermissions []PermissionApp `json:"application_permissions"`
}
