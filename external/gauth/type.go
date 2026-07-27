package gauth

type ClientSecret string

type ClientId string

type AuthRedirectUrl string

type Scopes string // the value should be json stringify array of string
type GenerateGauthConfigInput struct {
	ClientId
	ClientSecret
	AuthRedirectUrl
	*Scopes
}
