package jwt_handler

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

type CostumClaimsPayload struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}
