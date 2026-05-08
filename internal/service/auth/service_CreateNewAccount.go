package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (s *AuthService) CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error) {

	ctx, err = s.AuthRepositoryWriter.Begin(ctx, nil)
	if err != nil {
		err = errors.Join(err, ErrBeginDbTx)
		return
	}

	defer func() {
		err = s.AuthRepositoryWriter.CommitOrRollback(ctx, err)
	}()

	token, err := s.CreateToken(ctx, UserTokenClaimInput{
		TokenType:     TokenTypeJwt,
		UserId:        uuid.NewString(),
		ExpiredAtUtc0: time.Now().Add(time.Duration(s.JwtExpiration) * time.Minute).UnixMilli(),
	})
	if err != nil {
		return
	}
	_, err = s.AuthRepositoryWriter.CreateNewAccount(ctx, input)
	if err != nil {
		return
	}
	output.Id = token.Token
	return
}
