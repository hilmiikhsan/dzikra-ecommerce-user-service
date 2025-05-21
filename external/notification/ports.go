package notification

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/proto/notification"
)

type ExternalNotification interface {
	SendNotification(ctx context.Context, recipient, templateName string, placeholder map[string]string) error
	GetNotificationByType(ctx context.Context, req *notification.GetNotificationByTypeRequest) (*notification.GetNotificationByTypeResponse, error)
	CreateNotification(ctx context.Context, req *notification.CreateNotificationRequest) (*notification.CreateNotificationResponse, error)
	GetListNotification(ctx context.Context, req *notification.GetListNotificationRequest) (*notification.GetListNotificationResponse, error)
	SendFcmBatchNotification(ctx context.Context, req *notification.SendFcmBatchNotificationRequest) (*notification.SendFcmBatchNotificationResponse, error)
	SendFcmNotification(ctx context.Context, req *notification.SendFcmNotificationRequest) (*notification.SendFcmNotificationResponse, error)
}
