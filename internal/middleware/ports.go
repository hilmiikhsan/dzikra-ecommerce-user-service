package middleware

import "github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/jwt_handler"

type AuthMiddleware struct {
	jwt jwt_handler.JWT
}

func NewAuthMiddleware(jwt jwt_handler.JWT) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}
