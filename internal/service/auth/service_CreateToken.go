package auth

import "context"

func (s *AuthService) CreateToken(ctx context.Context, input UserTokenClaimInput) (output UserTokenClaimOutput, err error) {
	if input.TokenType == TokenTypeJwt {
		jwt, err := s.CreateJWTToken(ctx, CreateJWTTokenInput{
			ExpiredAtUtc0: input.ExpiredAtUtc0,
			IssuerAtUtc0:  input.IssuerAtUtc0,
			UserId:        input.UserId,
		})
		if err != nil {
			return UserTokenClaimOutput{}, err
		}
		output.Token = jwt.Token
	}
	return
}
