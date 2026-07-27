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

type GetAccessTokenByIdInput struct {
	AccessToken string //access token id
}

type GetAccessTokenByIdOutput struct {
	Id           string
	Username     *string
	Email        *string
	AccessToken  string
	RefreshToken string
	ExpiryAtUtc0 *int64
}
