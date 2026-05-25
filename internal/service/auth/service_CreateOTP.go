package auth

import (
	"context"
	"time"

	"github.com/masraga/kerp-api/internal/util/generator"
)

func (s *AuthService) CreateOTP(ctx context.Context, input CreateOTPInput) (output CreateOTPOutput, err error) {
	authUser, err := s.AuthRepositoryReader.FindAuth(ctx, FindAuthInput{
		UserId: input.UserId,
	})
	if err != nil {
		return
	}
	if authUser.Id == "" {
		err = ErrAuthNotFound
		return
	}

	_, err = s.AuthRepositoryWriter.DeleteAllUserOTP(ctx, DeleteAllUserOTPInput{
		UserId: input.UserId,
	})
	if err != nil {
		return
	}

	otp := input.OtpCode
	if otp == "" {
		otp, err = generator.GenerateRandom(6, true)
		if err != nil {
			return
		}
	}

	expiredAtUtc0 := input.ExpiredAtUtc0
	if expiredAtUtc0 == 0 {
		expiredAtUtc0 = time.Now().Add(time.Duration(OtpExpiredDuration) * time.Minute).UnixMilli()
	}

	_, err = s.AuthRepositoryWriter.CreateOTP(ctx, CreateOTPInput{
		OtpCode:       otp,
		UserId:        input.UserId,
		Note:          input.Note,
		ExpiredAtUtc0: expiredAtUtc0,
	})
	if err != nil {
		return
	}

	output.OtpCode = otp

	return
}
