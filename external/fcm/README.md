# FCM Service

Firebase Cloud Messaging is the current push provider for the notification service.

## Rules

- `FcmService` must implement [`PushProviderInterface`](../../internal/service/notification/interface.go).
- Application code must call `NotificationService.PushNotification` instead of depending directly on this package.
- `SendNotificationInput.UserId` is used as the FCM device registration token.
- `FCM_SERVICE_ACCOUNT_ID` supplies the Firebase project ID.
- Firebase SDK default credentials must be available in the runtime environment.
- Provider errors must include the notification domain errors defined in [`error.go`](../../internal/service/notification/error.go).

See [Push Notifications](../../.codex/tech-docs/features/push-notifications.md) for the full service flow and current payload behavior.
