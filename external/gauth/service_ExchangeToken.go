package gauth

import (
	"context"

	"github.com/google/uuid"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"golang.org/x/oauth2"
)

func (s *GAuthService) ExchangeToken(ctx context.Context, input oauth.GoogleOauthCodeExchangeInput) (output oauth.GoogleOauthCodeExchangeOutput, err error) {
	oauthConfig, err := GenerateGauthConfig(GenerateGauthConfigInput{
		ClientId:        s.clientId,
		ClientSecret:    s.clientSecret,
		AuthRedirectUrl: s.authRedirectUrl,
		Scopes:          &s.scopes,
	})
	if err != nil {
		err = s.err.Wrap(err)
		return
	}
	state := uuid.New().String()

	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	output.Url = url
	return
}
