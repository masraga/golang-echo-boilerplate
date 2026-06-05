package notification

import "context"

type NotificationServiceInterface interface {
	PushNotification(ctx context.Context, input SendNotificationInput) (output SendNotificationOutput, err error)
}

type PushProviderInterface interface {
	SendNotification(ctx context.Context, input SendNotificationInput) (output SendNotificationOutput, err error)
}
