package auth

import "context"

func (s *AuthService) StoreAccessToken(ctx context.Context, input StoreAccessTokenInput) (output StoreAccessTokenOutput, err error) {
	output, err = s.AuthRepositoryWriter.StoreAccessToken(ctx, input)
	return
}
