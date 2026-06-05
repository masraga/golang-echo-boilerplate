package fcm_test

import (
	"context"
	"testing"

	"github.com/masraga/kerp-api/external/fcm"
	"github.com/masraga/kerp-api/internal/service/notification"
	"github.com/masraga/kerp-api/internal/testutil"
	"github.com/rs/zerolog"
)

func TestFcmService_SendNotification(t *testing.T) {
	var (
		expectedProjectId string = "kerp-api-f5384"
		expectedUserId    string = "dwmD6py65-COEvC8xMEsjo:APA91bHQClUnYKkSpZyOG1p6Ya2uVtM0CDO8pkxK7eAyOBQKe9oNJAAIzypbEQ4-78LjipN5QOAUxEawgoy0FJYr-Jglu79qQXLpbsMxxefInT_wjw5pP4w"
	)

	type args struct {
		ctx   context.Context
		input notification.SendNotificationInput
	}

	type expected = testutil.Result[notification.SendNotificationOutput]

	type test struct {
		name     string
		args     args
		expected expected
	}

	tests := []test{
		{
			name: "successfully send message",
			args: args{
				ctx: context.Background(),
				input: notification.SendNotificationInput{
					UserId: expectedUserId,
					Title:  "title notif",
					Body:   "body notif",
				},
			},
			expected: expected{
				Err:   nil,
				Value: notification.SendNotificationOutput{IsSuccess: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := fcm.NewFcmService(fcm.FcmServiceStructOpts{
				Ctx:              tt.args.ctx,
				ServiceAccountId: fcm.ConfigServiceAccountId(expectedProjectId),
				Logger:           zerolog.Nop(),
			})

			got, err := svc.SendNotification(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
