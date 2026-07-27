package oauth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func (s *OAuthService) Callback(ctx context.Context, input GoogleOauthCallbackInput) (output GoogleOauthCallbackOutput, err error) {
	providerOutput, err := s.provider.Callback(ctx, input)
	if err != nil {
		err = s.err.Wrap(errors.Join(err, ErrFailedGetCallbackToken))
		return
	}

	id := uuid.New().String()
	_, err = s.oauthRepositoryWriter.CreateAccessToken(ctx, CreateAccessTokenInput{
		Id:           id,
		AccessToken:  providerOutput.Token,
		RefreshToken: providerOutput.RefreshToken,
	})

	if err != nil {
		err = s.err.Wrap(err)
		return
	}

	output.Token = providerOutput.Token
	output.RefreshToken = providerOutput.RefreshToken

	return
}
