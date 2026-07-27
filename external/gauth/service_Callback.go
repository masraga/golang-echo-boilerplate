package gauth

import (
	"context"

	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
)

func (s *GAuthService) Callback(ctx context.Context, input oauth.GoogleOauthCallbackInput) (output oauth.GoogleOauthCallbackOutput, err error) {
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

	token, err := oauthConfig.Exchange(ctx, input.Code)
	if err != nil {
		return
	}

	output.Token = token.AccessToken
	output.RefreshToken = token.RefreshToken

	return
}
