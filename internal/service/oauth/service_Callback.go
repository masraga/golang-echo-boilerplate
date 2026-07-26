package oauth

import "context"

func (s *OAuthService) Callback(ctx context.Context, input GoogleOauthCallbackInput) (output GoogleOauthCallbackOutput, err error) {
	providerOutput, err := s.provider.Callback(ctx, input)
	if err != nil {
		return
	}

	output.Token = providerOutput.Token
	output.RefreshToken = providerOutput.RefreshToken

	return
}
