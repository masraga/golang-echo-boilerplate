package notification

import "errors"

var (
	ErrInitiateNotifProvider error = errors.New("error when initiate notification provider")
	ErrToSendNotification    error = errors.New("error when send notification")
)
