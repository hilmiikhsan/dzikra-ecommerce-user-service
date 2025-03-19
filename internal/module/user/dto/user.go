package dto

import (
	role "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/dto"
)

type GetListUserResponse struct {
	Users       []GetListUser `json:"users"`
	TotalPages  int           `json:"total_pages"`
	CurrentPage int           `json:"current_page"`
	PageSize    int           `json:"page_size"`
	TotalData   int           `json:"total_data"`
}

type GetListUser struct {
	ID             string          `json:"id"`
	Email          string          `json:"email"`
	FullName       string          `json:"full_name"`
	PhoneNumber    string          `json:"phone_number"`
	UserRole       []role.UserRole `json:"user_role"`
	EmailConfirmed IsConfirmEmail  `json:"email_confirmed"`
}
