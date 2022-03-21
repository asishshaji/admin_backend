package notification_service

import (
	"context"

	"github.com/asishshaji/admin-api/models"
)

type INotificationService interface {
	SendNotification(ctx context.Context, message models.NotificationMessage) error
}
