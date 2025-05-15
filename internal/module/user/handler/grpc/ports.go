package grpc

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/tokenvalidation"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type userGrpcAPI struct {
	jwt jwt_handler.JWT
	tokenvalidation.UnimplementedTokenValidationServer
}

func NewUserGrpcAPI() *userGrpcAPI {
	var handler = new(userGrpcAPI)

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwt_handler.NewJWT(redisRepository)

	// grpc handler
	handler.jwt = jwt

	return handler
}
