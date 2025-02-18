package middleware

import "github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/jwt_handler"

type UserMiddleware struct {
	jwt jwt_handler.JWT
}

func NewUserMiddleware(jwt jwt_handler.JWT) *UserMiddleware {
	return &UserMiddleware{
		jwt: jwt,
	}
}
