package notification_service

import (
	"context"
	"log"
	"os"

	"github.com/asishshaji/admin-api/models"
	onesignal "github.com/tbalthazar/onesignal-go"
)

type NotificationService struct {
	l      *log.Logger
	client *onesignal.Client
	appId  string
}

func createClient() *onesignal.Client {
	client := onesignal.NewClient(nil)
	client.AppKey = os.Getenv("ONE_SIGNAL_KEY")

	return client
}

func NewNotificationService(l *log.Logger) INotificationService {
	return NotificationService{
		l:      l,
		client: createClient(),
		appId:  os.Getenv("ONE_SIGNAL_APP_ID"),
	}
}

func (nS NotificationService) SendNotification(ctx context.Context, msg models.NotificationMessage) error {
	notificationReq := &onesignal.NotificationRequest{
		AppID:            nS.appId,
		Contents:         msg.Contents,
		Headings:         msg.Heading,
		IncludePlayerIDs: []string{msg.UserToken},
	}

	createRes, res, err := nS.client.Notifications.Create(notificationReq)

	if err != nil {
		nS.l.Println(err)
		return err
	}
	nS.l.Println(createRes, res)

	return nil
}
