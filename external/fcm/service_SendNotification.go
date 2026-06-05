package fcm

import (
	"context"
	"errors"

	"firebase.google.com/go/v4/messaging"
	"github.com/masraga/kerp-api/internal/service/notification"
)

func (s *FcmService) SendNotification(ctx context.Context, input notification.SendNotificationInput) (output notification.SendNotificationOutput, err error) {
	msg := &messaging.Message{
		Token: input.UserId,
		Notification: &messaging.Notification{
			Title: input.Title,
			Body:  input.Body,
		},
		Data: map[string]string{
			"type": "general",
		},
	}
	_, err = s.client.Send(ctx, msg)
	if err != nil {
		err = errors.Join(err, notification.ErrToSendNotification)
		return
	}
	output.IsSuccess = true
	return
}
