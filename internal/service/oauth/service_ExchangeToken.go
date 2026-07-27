package oauth

import (
	"context"
	"errors"
)

func (s *OAuthService) ExchangeToken(ctx context.Context, input GoogleOauthCodeExchangeInput) (output GoogleOauthCodeExchangeOutput, err error) {
	providerOutput, err := s.provider.ExchangeToken(ctx, input)
	if err != nil {
		err = s.err.Wrap(errors.Join(err, ErrFailedExchangeToken))
		return
	}
	output.Url = providerOutput.Url
	return
}
