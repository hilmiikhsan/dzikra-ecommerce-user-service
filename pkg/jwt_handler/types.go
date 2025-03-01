package jwt_handler

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	jwt.RegisteredClaims
}

type CostumClaimsPayload struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type GenerateTokenResponse struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	TokenExpiredAt        time.Time `json:"expired_at"`
	RefreshTokenExpiredAt time.Time `json:"expired_at"`
	CreatedAt             time.Time `json:"created_at"`
}
