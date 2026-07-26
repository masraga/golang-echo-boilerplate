package gauth

import (
	"context"

	"github.com/google/uuid"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (s *GAuthService) ExchangeToken(ctx context.Context, input oauth.GoogleOauthCodeExchangeInput) (output oauth.GoogleOauthCodeExchangeOutput, err error) {
	var oauthConfig = &oauth2.Config{
		ClientID:     string(s.clientId),
		ClientSecret: string(s.clientSecret),
		RedirectURL:  string(s.authRedirectUrl),
		Scopes: []string{
			"https://www.googleapis.com/auth/drive.readonly",
			"https://www.googleapis.com/auth/spreadsheets",
		},
		Endpoint: google.Endpoint,
	}
	state := uuid.New().String()

	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	output.Url = url
	return
}
