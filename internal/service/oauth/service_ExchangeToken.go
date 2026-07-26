package oauth

import "context"

func (s *OAuthService) ExchangeToken(ctx context.Context, input GoogleOauthCodeExchangeInput) (output GoogleOauthCodeExchangeOutput, err error) {
	providerOutput, err := s.provider.ExchangeToken(ctx, input)
	if err != nil {
		return
	}
	output.Url = providerOutput.Url
	return
}
