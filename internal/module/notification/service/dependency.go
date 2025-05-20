package service

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/notification"
	notificationPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/notification/ports"
	userFcmTokenPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user_fcm_token/ports"
	"github.com/jmoiron/sqlx"
)

var _ notificationPorts.NotificationService = &notificationService{}

type notificationService struct {
	db                     *sqlx.DB
	externalNotification   externalNotification.ExternalNotification
	userFcmTokenRepository userFcmTokenPorts.UserFCMTokenRepository
	// fcmClient              *messaging.Client
}

func NewNotificationService(
	db *sqlx.DB,
	externalNotification externalNotification.ExternalNotification,
	userFcmTokenRepository userFcmTokenPorts.UserFCMTokenRepository,
	// fcmClient *messaging.Client,
) *notificationService {
	return &notificationService{
		db:                     db,
		externalNotification:   externalNotification,
		userFcmTokenRepository: userFcmTokenRepository,
		// fcmClient:              fcmClient,
	}
}
