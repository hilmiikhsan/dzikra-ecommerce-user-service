package service

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/notification/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *notificationService) CreateNotification(ctx context.Context, req *dto.CreateNotificationRequest) error {
	_, err := s.externalNotification.CreateNotification(ctx, &notification.CreateNotificationRequest{
		Title:   req.Title,
		Detail:  req.Detail,
		Url:     req.Url,
		NTypeId: req.NTypeID,
		UserId:  req.UserID,
	})
	if err != nil {
		log.Error().Err(err).Msg("service::CreateNotification - Failed to create notification")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}

func (s *notificationService) GetListNotification(ctx context.Context, page, limit int, search string) (*dto.GetListNotificationResponse, error) {
	res, err := s.externalNotification.GetListNotification(ctx, &notification.GetListNotificationRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Search: search,
	})
	if err != nil {
		log.Error().Err(err).Msg("service::GetListNotification - Failed to get list notification")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	notifications := make([]dto.NotificationDetail, 0, len(res.Notifications))
	for _, notification := range res.Notifications {
		var detail *string
		if notification.Detail != "" {
			detail = &notification.Detail
		} else {
			detail = nil
		}

		notifications = append(notifications, dto.NotificationDetail{
			ID:        int(notification.Id),
			Title:     notification.Title,
			Detail:    detail,
			Url:       notification.Url,
			NTypeID:   notification.NTypeId,
			UserID:    notification.UserId,
			CreatedAt: notification.CreatedAt,
		})
	}

	return &dto.GetListNotificationResponse{
		Notification: notifications,
		TotalPages:   int(res.TotalPage),
		CurrentPage:  int(res.CurrentPage),
		PageSize:     int(res.PageSize),
		TotalData:    int(res.TotalData),
	}, nil
}

func (s *notificationService) SendFcmBatchNotification(ctx context.Context, req *dto.SendFcmBatchRequest) error {
	role := constants.RoleSuperAdminPOS
	if req.IsUser {
		role = constants.UserRole
	}

	fcmTokens, err := s.userFcmTokenRepository.FindFcmUserTokenByRole(ctx, role)
	if err != nil {
		log.Error().Err(err).Msg("service::SendFcmBatchNotification - Failed to find fcm user token by role")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if len(fcmTokens) == 0 {
		log.Warn().Msg("service::SendFcmBatchNotification - No fcm user token found")
		return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrUserNotFound))
	}

	go func() {
		_, err = s.externalNotification.SendFcmBatchNotification(ctx, &notification.SendFcmBatchNotificationRequest{
			FcmTokens: fcmTokens,
			Title:     req.Title,
			Body:      req.Detail,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::SendFcmBatchNotification - Failed to send fcm batch notification")
		}
	}()

	log.Info().Msg("service::SendFcmBatchNotification - Success to send fcm batch notification")

	return nil
}
