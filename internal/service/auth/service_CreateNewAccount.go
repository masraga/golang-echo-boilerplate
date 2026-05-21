package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (s *AuthService) CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error) {

	authUser, _ := s.AuthRepositoryReader.FindAuth(ctx, FindAuthInput{
		PhoneNo: input.PhoneNo,
	})
	if authUser.Id != "" {
		err = ErrDuplicateUser
		return
	}

	ctx, err = s.AuthRepositoryWriter.Begin(ctx, nil)
	if err != nil {
		err = errors.Join(err, ErrBeginDbTx)
		return
	}

	defer func() {
		commitErr := s.AuthRepositoryWriter.CommitOrRollback(ctx, err)
		err = commitErr
	}()

	input.Id = uuid.NewString()

	_, err = s.AuthRepositoryWriter.CreateNewAccount(ctx, input)
	if err != nil {
		return
	}

	token, err := s.CreateToken(ctx, UserTokenClaimInput{
		TokenType:     TokenTypeJwt,
		UserId:        input.Id,
		ExpiredAtUtc0: time.Now().Add(time.Duration(s.JwtExpiration) * time.Minute).UnixMilli(),
	})
	if err != nil {
		return
	}

	output.Id = token.Token
	return
}
