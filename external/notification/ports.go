package external

import "context"

type ExternalNotification interface {
	SendNotification(ctx context.Context, recipient, templateName string, placeholder map[string]string) error
}
