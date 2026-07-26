package gauth

import (
	"context"

	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (s *GAuthService) Callback(ctx context.Context, input oauth.GoogleOauthCallbackInput) (output oauth.GoogleOauthCallbackOutput, err error) {
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

	token, err := oauthConfig.Exchange(ctx, input.Code)
	if err != nil {
		return
	}

	output.Token = token.AccessToken
	output.RefreshToken = token.RefreshToken

	return
}
