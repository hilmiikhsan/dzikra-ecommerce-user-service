package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/notification/dto"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, req *dto.CreateNotificationRequest) error
	GetListNotification(ctx context.Context, page, limit int, search string) (*dto.GetListNotificationResponse, error)
	SendFcmBatchNotification(ctx context.Context, req *dto.SendFcmBatchRequest) error
}
