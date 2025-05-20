package rest

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/redis"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/notification/ports"
	notificationService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/notification/service"
	userFcmTokenRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user_fcm_token/repository"
	jwtHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/jwt_handler"
)

type notificationHandler struct {
	service    ports.NotificationService
	middleware middleware.UserMiddleware
	validator  adapter.Validator
}

func NewNotificationHandler() *notificationHandler {
	var handler = new(notificationHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalOrder := &externalNotification.External{}

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewUserMiddleware(jwt)

	// firebase messaging
	// fcmClient := firebase.InitFirebaseMessaging()

	// repository
	userFcmTokenRepository := userFcmTokenRepository.NewUserFcmTokenRepository(adapter.Adapters.DzikraPostgres)

	// notification  service
	notificationService := notificationService.NewNotificationService(
		adapter.Adapters.DzikraPostgres,
		externalOrder,
		userFcmTokenRepository,
		// fcmClient,
	)

	// handler
	handler.service = notificationService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
