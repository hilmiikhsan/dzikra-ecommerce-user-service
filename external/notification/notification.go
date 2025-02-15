package external

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/external/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
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
