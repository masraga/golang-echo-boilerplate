package notification

import "context"

func (s *NotificationService) PushNotification(ctx context.Context, input SendNotificationInput) (output SendNotificationOutput, err error) {
	_, err = s.Provider.SendNotification(ctx, input)
	output.IsSuccess = err == nil
	return
}
