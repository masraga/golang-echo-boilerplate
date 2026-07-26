package oauth

type GoogleOauthCodeExchangeInput struct {
}

type GoogleOauthCodeExchangeOutput struct {
	Url string
}

type GoogleOauthCallbackInput struct {
	Code string
}

type GoogleOauthCallbackOutput struct {
	Token        string
	RefreshToken string
}
