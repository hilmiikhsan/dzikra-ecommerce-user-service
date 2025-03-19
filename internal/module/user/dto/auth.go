package dto

import (
	role "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/dto"
)

// User DTO
type RegisterRequest struct {
	FullName        string `json:"full_name" validate:"required,min=2,max=100"`
	Username        string `json:"username" validate:"required,min=2,max=50"`
	Password        string `json:"password" validate:"required,strong_password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email,email_blacklist"`
	PhoneNumber     string `json:"phone_number" validate:"required,phone"`
}

type RegisterResponse struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

type VerificationRequest struct {
	Email string `json:"email" validate:"required,email,email_blacklist"`
	Otp   string `json:"otp" validate:"required,otp_number"`
}

type VerificationResponse struct {
	Email          string         `json:"email"`
	EmailConfirmed EmailConfirmed `json:"email_confirmed"`
}

type EmailConfirmed struct {
	IsConfirm bool   `json:"is_confirm"`
	CreatedAt string `json:"created_at"`
}

type SendOtpNumberVerificationRequest struct {
	Email string `json:"email" validate:"required,email,email_blacklist"`
}

type SendOtpNumberVerificationResponse struct {
	Otp string `json:"otp"`
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	DeviceID   string `json:"device_id" validate:"required,max=100"`
	DeviceType string `json:"device_type" validate:"required,max=10,device_type"`
	FcmToken   string `json:"fcm_token" validate:"required,max=255"`
}

type AuthUserResponse struct {
	Email          string                `json:"email"`
	EmailConfirmed EmailConfirmed        `json:"email_confirmed"`
	FullName       string                `json:"full_name"`
	PhoneNumber    string                `json:"phone_number"`
	Token          TokenDetail           `json:"token"`
	UserRole       []role.UserRoleDetail `json:"user_role"`
}

type TokenDetail struct {
	Token        string             `json:"token"`
	ExpiredAt    string             `json:"expired_at"`
	CreatedAt    string             `json:"created_at"`
	RefreshToken RefreshTokenDetail `json:"refresh_token"`
}

type RefreshTokenDetail struct {
	RefreshToken string `json:"refresh_token"`
	ExpiredAt    string `json:"expired_at"`
	CreatedAt    string `json:"created_at"`
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

type GetCurrentUserResponse struct {
	Email          string                `json:"email"`
	EmailConfirmed EmailConfirmed        `json:"email_confirmed"`
	FullName       string                `json:"full_name"`
	PhoneNumber    string                `json:"phone_number"`
	UserRole       []role.UserRoleDetail `json:"user_role"`
}

type ForgotPasswordResponse struct {
	Email    string `json:"email"`
	Sessions string `json:"sessions"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" validate:"required,email"`
	SessionToken    string `json:"session_token" validate:"required,max=255"`
	Password        string `json:"password" validate:"required,strong_password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type IsConfirmEmail struct {
	IsConfirm bool `json:"is_confirm"`
}
