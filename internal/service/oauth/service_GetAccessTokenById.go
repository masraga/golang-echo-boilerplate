package oauth

import "context"

func (s *OAuthService) GetAccessTokenById(ctx context.Context, input GetAccessTokenByIdInput) (output GetAccessTokenByIdOutput, err error) {
	output, err = s.oauthRepositoryReader.GetAccessTokenById(ctx, input)
	if err != nil {
		err = s.err.Wrap(err)
		return
	}

	return
}
