# Push Notifications

Push notifications use a provider-independent notification service with Firebase Cloud Messaging as the currently wired provider. This is an internal application capability and does not expose a public HTTP endpoint.

## Architecture

| Layer | File | Purpose |
| --- | --- | --- |
| Service contract | `internal/service/notification/interface.go` | Exposes `NotificationServiceInterface.PushNotification` and the provider boundary `PushProviderInterface.SendNotification`. |
| Service | `internal/service/notification/service_PushNotification.go` | Delegates the input to the configured provider and sets `IsSuccess` from the returned error. |
| Provider | `external/fcm/service_SendNotification.go` | Converts the shared input into a Firebase message and sends it with the Firebase Messaging client. |
| Provider initialization | `external/fcm/service.go` | Creates the Firebase app and Messaging client for the configured project. |
| Dependency injection | `app/backend/wire_provider.go` | Constructs FCM as the application's `PushProviderInterface`. |

Application code must call `NotificationServiceInterface.PushNotification`. It must not depend directly on `external/fcm`; provider-specific behavior stays behind `PushProviderInterface`.

## Input And Delivery

`SendNotificationInput` contains `UserId`, `Title`, `Body`, and optional `Icon`.

- `UserId` is currently used as the FCM device registration token.
- `Title` and `Body` populate the Firebase notification payload.
- FCM sends a data field with `type=general`.
- `Icon` is part of the shared input but is not currently mapped by the FCM provider.

`PushNotification` returns the provider error unchanged. `SendNotificationOutput.IsSuccess` is `true` only when the provider returns no error.

The FCM provider joins delivery failures with `ErrToSendNotification`. Provider initialization failures are logged with `ErrInitiateNotifProvider`, and `NewFcmService` returns `nil`.

## Configuration

The backend reads `FCM_SERVICE_ACCOUNT_ID` and passes it to Firebase as the project ID. Because no explicit credential option is supplied, Firebase SDK default credential discovery must be configured in the runtime environment.

The provider is initialized during backend dependency injection. A failed initialization produces a nil provider, so valid Firebase configuration and credentials are required before the notification service is used.

## Tests

`external/fcm/service_SendNotification_test.go` exercises a real FCM send. It requires valid Firebase credentials, network access, the configured project, and a usable device registration token.

When changing `NotificationServiceInterface` or `PushProviderInterface`, regenerate mocks with:

```sh
make clean init
```
