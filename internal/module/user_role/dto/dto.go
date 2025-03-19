package dto

type UserRole struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

type UserRoleDetail struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}
