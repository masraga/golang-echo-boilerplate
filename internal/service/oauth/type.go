package oauth

import "github.com/masraga/golang-echo-boilerplate/internal/service/auth"

type ExchangeTokenAction string
type GoogleOauthCodeExchangeInput struct {
	Actions []ExchangeTokenAction
}

type GoogleOauthCodeExchangeOutput struct {
	Url string
}

type GoogleOauthCallbackInput struct {
	Code  string
	State string
}

type GoogleOauthCallbackOutput struct {
	IdToken      string // used for get user profile data
	Email        string `json:"email"`
	Name         string `json:"name"`
	Sub          string `json:"sub"`
	Picture      string `json:"picture"`
	Token        string
	RefreshToken string
	State        string
}

type GoogleOauthCallbackState struct {
	Actions []ExchangeTokenAction
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

type CreateAccessTokenInput struct {
	Id           string
	AccessToken  string
	RefreshToken string
}

type CreateAccessTokenOutput struct {
	Id string
}

type BypassCreateNewUserInput struct {
	Id        string
	PhoneNo   string
	Pin       string
	Email     string
	CreatedBy string
	LoginMode auth.LoginMode
}

type BypassCreateNewUserOutput struct {
	Id string
}
