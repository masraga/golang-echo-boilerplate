package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/masraga/kerp-api/internal/util/generator"
	"github.com/masraga/kerp-api/internal/util/pointer"
)

func (s *AuthService) CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error) {

	authUser, err := s.AuthRepositoryReader.FindAuth(ctx, FindAuthInput{
		PhoneNo: input.PhoneNo,
	})
	if err != nil && !errors.Is(err, ErrAuthNotFound) {
		return
	}
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

	otpCode, err := generator.GenerateRandom(6, true)
	if err != nil {
		return
	}

	otpSvc, err := s.CreateOTP(ctx, CreateOTPInput{
		UserId:        input.Id,
		Note:          pointer.String("OTP for new account register"),
		ExpiredAtUtc0: time.Now().Add(time.Duration(OtpExpiredDuration) * time.Minute).UnixMilli(),
		OtpCode:       otpCode,
	})
	if err != nil {
		return
	}

	output.Id = input.Id
	output.OtpCode = otpSvc.OtpCode
	return
}
