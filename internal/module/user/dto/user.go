package dto

import (
	role "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/dto"
)

type GetListUserResponse struct {
	Users       []GetDetailUserResponse `json:"users"`
	TotalPages  int                     `json:"total_pages"`
	CurrentPage int                     `json:"current_page"`
	PageSize    int                     `json:"page_size"`
	TotalData   int                     `json:"total_data"`
}

type GetDetailUserResponse struct {
	ID             string          `json:"id"`
	Email          string          `json:"email"`
	FullName       string          `json:"full_name"`
	PhoneNumber    string          `json:"phone_number"`
	UserRole       []role.UserRole `json:"user_role"`
	EmailConfirmed IsConfirmEmail  `json:"email_confirmed"`
}

type CreateUserRequest struct {
	FullName        string   `json:"full_name" validate:"required,min=2,max=100"`
	Email           string   `json:"email" validate:"required,email,email_blacklist"`
	Password        string   `json:"password" validate:"required,strong_password"`
	ConfirmPassword string   `json:"confirm_password" validate:"required,eqfield=Password"`
	Role            []string `json:"role" validate:"non_empty_array"`
	PhoneNumber     string   `json:"phone_number" validate:"required,phone"`
}

type CreateUserResponse struct {
	ID          string          `json:"id"`
	Email       string          `json:"email"`
	FullName    string          `json:"full_name"`
	PhoneNumber string          `json:"phone_number"`
	UserRole    []role.UserRole `json:"user_role"`
}
