package oauth

import (
	"context"
	"errors"
)

func (s *OAuthService) Callback(ctx context.Context, input GoogleOauthCallbackInput) (output GoogleOauthCallbackOutput, err error) {
	providerOutput, err := s.provider.Callback(ctx, input)
	if err != nil {
		err = s.err.Wrap(errors.Join(err, ErrFailedGetCallbackToken))
		return
	}

	output.Token = providerOutput.Token
	output.RefreshToken = providerOutput.RefreshToken

	return
}
