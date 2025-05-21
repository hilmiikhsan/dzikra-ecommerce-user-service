package notification

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) SendNotification(ctx context.Context, recipient, templateName string, placeholder map[string]string) error {
	conn, err := grpc.Dial(utils.GetEnv("NOTIFICATION_GRPC_HOST", config.Envs.Notification.NotificationGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::SendNotification - Failed to dial grpc")
		return err
	}
	defer conn.Close()

	client := notification.NewNotificationServiceClient(conn)
	request := &notification.SendNotificationRequest{
		Recipient:    recipient,
		TemplateName: templateName,
		Placeholders: placeholder,
	}

	resp, err := client.SendNotification(ctx, request)
	if err != nil {
		log.Err(err).Msg("external::SendNotification - Failed to send notification")
		return err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::SendNotification - Response error from notification")
		return fmt.Errorf("get response error from notification: %s", resp.Message)
	}

	return nil
}

func (*External) GetNotificationByType(ctx context.Context, req *notification.GetNotificationByTypeRequest) (*notification.GetNotificationByTypeResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("NOTIFICATION_GRPC_HOST", config.Envs.Notification.NotificationGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::GetNotificationByType - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := notification.NewNotificationServiceClient(conn)
	resp, err := client.GetNotificationByType(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::GetNotificationByType - Failed to get notification by type")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		if resp.Message == constants.ErrNotificationTypeNotFound {
			log.Err(err).Msg("external::GetNotificationByType - Notification type not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrNotificationTypeNotFound))
		}

		log.Err(err).Msg("external::GetNotificationByType - Response error from notification")
		return nil, fmt.Errorf("get response error from notification: %s", resp.Message)
	}

	return resp, nil
}

func (*External) CreateNotification(ctx context.Context, req *notification.CreateNotificationRequest) (*notification.CreateNotificationResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("NOTIFICATION_GRPC_HOST", config.Envs.Notification.NotificationGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::CreateNotification - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := notification.NewNotificationServiceClient(conn)
	resp, err := client.CreateNotification(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::CreateNotification - Failed to create notification")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		if resp.Message == constants.ErrNotificationTypeNotFound {
			log.Err(err).Msg("external::GetNotificationByType - Notification type not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrNotificationTypeNotFound))
		}

		log.Err(err).Msg("external::CreateNotification - Response error from notification")
		return nil, fmt.Errorf("get response error from notification: %s", resp.Message)
	}

	return resp, nil
}

func (*External) GetListNotification(ctx context.Context, req *notification.GetListNotificationRequest) (*notification.GetListNotificationResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("NOTIFICATION_GRPC_HOST", config.Envs.Notification.NotificationGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::GetListNotification - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := notification.NewNotificationServiceClient(conn)
	resp, err := client.GetListNotification(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::GetListNotification - Failed to get list notification")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::GetListNotification - Response error from notification")
		return nil, fmt.Errorf("get response error from notification: %s", resp.Message)
	}

	return resp, nil
}

func (*External) SendFcmBatchNotification(ctx context.Context, req *notification.SendFcmBatchNotificationRequest) (*notification.SendFcmBatchNotificationResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("NOTIFICATION_GRPC_HOST", config.Envs.Notification.NotificationGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::SendFcmBatchNotification - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := notification.NewNotificationServiceClient(conn)
	resp, err := client.SendFcmBatchNotification(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::SendFcmBatchNotification - Failed to send fcm batch notification")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::SendFcmBatchNotification - Response error from notification")
		return nil, fmt.Errorf("get response error from notification: %s", resp.Message)
	}

	return resp, nil
}

func (*External) SendFcmNotification(ctx context.Context, req *notification.SendFcmNotificationRequest) (*notification.SendFcmNotificationResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("NOTIFICATION_GRPC_HOST", config.Envs.Notification.NotificationGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::SendFcmNotification - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := notification.NewNotificationServiceClient(conn)
	resp, err := client.SendFcmNotification(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::SendFcmNotification - Failed to send fcm notification")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::SendFcmNotification - Response error from notification")
		return nil, fmt.Errorf("get response error from notification: %s", resp.Message)
	}

	return resp, nil
}
