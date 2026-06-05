package notification

type NotificationService struct {
	Provider PushProviderInterface
}

type NotificationServiceOpts struct {
	Provider PushProviderInterface
}

func NewNotificationService(opts NotificationServiceOpts) *NotificationService {
	return &NotificationService{
		Provider: opts.Provider,
	}
}
