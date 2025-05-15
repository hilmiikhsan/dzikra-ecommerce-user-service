package jwt_handler

import (
	"time"

	role "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role/dto"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID     string                `json:"user_id"`
	Email      string                `json:"email"`
	FullName   string                `json:"full_name"`
	SessionID  string                `json:"session_id"`
	DeviceID   string                `json:"device_id"`
	DeviceType string                `json:"device_type"`
	FcmToken   string                `json:"fcm_token"`
	UserRoles  []role.UserRoleDetail `json:"user_roles"`
	CreatedAt  time.Time             `json:"created_at"`
	jwt.RegisteredClaims
}

type CostumClaimsPayload struct {
	UserID     string                `json:"user_id"`
	Email      string                `json:"email"`
	FullName   string                `json:"full_name"`
	SessionID  string                `json:"session_id"`
	DeviceID   string                `json:"device_id"`
	DeviceType string                `json:"device_type"`
	FcmToken   string                `json:"fcm_token"`
	UserRoles  []role.UserRoleDetail `json:"user_roles"`
}

type GenerateTokenResponse struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	TokenExpiredAt        time.Time `json:"expired_at"`
	RefreshTokenExpiredAt time.Time `json:"refresh_expired_at"`
	SessionID             string    `json:"session_id"`
	DeviceID              string    `json:"device_id"`
	DeviceType            string    `json:"device_type"`
	FcmToken              string    `json:"fcm_token"`
	CreatedAt             time.Time `json:"created_at"`
}
