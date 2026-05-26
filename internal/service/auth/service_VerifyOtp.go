package auth

import (
	"context"
	"time"
)

func (s *AuthService) VerifyOtp(ctx context.Context, input VerifyOtpInput) (output VerifyOtpOutput, err error) {

	authUser, err := s.AuthRepositoryReader.FindAuth(ctx, FindAuthInput{
		PhoneNo: input.PhoneNo,
	})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	otpUser, err := s.AuthRepositoryReader.FindOTP(ctx, FindOTPInput{
		UserId:  authUser.Id,
		OtpCode: input.OtpCode,
	})
	if err != nil {
		return
	}
	now := time.Now().UnixMilli()
	if now > otpUser.ExpiredAtUtc0 {
		err = s.Err.Wrap(ErrOtpIsExpired)
		return
	}

	if otpUser.IsVerified {
		err = s.Err.Wrap(ErrOtpAleadyVerified)
		return
	}

	verify, err := s.AuthRepositoryWriter.VerifyOtp(ctx, VerifyOtpInput{
		UserId:  authUser.Id,
		OtpCode: input.OtpCode,
	})
	if err != nil {
		return
	}

	output.UserId = verify.UserId
	output.IsValid = verify.IsValid
	output.Note = otpUser.Note
	return
}
