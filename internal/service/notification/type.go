package notification

type SendNotificationInput struct {
	UserId string
	Title  string
	Body   string
	Icon   *string
}

type SendNotificationOutput struct {
	IsSuccess bool
}
