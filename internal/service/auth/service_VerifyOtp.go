package auth

import (
	"context"
	"errors"
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

	_, err = s.AuthRepositoryWriter.UpdateOtpValidity(ctx, UpdateOtpValidityInput{
		UserId:     authUser.Id,
		IsOtpValid: false,
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

	ctx, err = s.AuthRepositoryWriter.Begin(ctx, nil)
	if err != nil {
		err = s.Err.Wrap(errors.Join(err, ErrBeginDbTx))
		return
	}

	defer func() {
		commitErr := s.AuthRepositoryWriter.CommitOrRollback(ctx, err)
		err = s.Err.Wrap(commitErr)
	}()

	verify, err := s.AuthRepositoryWriter.VerifyOtp(ctx, VerifyOtpInput{
		UserId:  authUser.Id,
		OtpCode: input.OtpCode,
	})
	if err != nil {
		return
	}

	_, err = s.AuthRepositoryWriter.VerifyUserAccount(ctx, VerifyUserAccountInput{
		PhoneNo: authUser.PhoneNo,
	})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	_, err = s.AuthRepositoryWriter.UpdateOtpValidity(ctx, UpdateOtpValidityInput{
		UserId:     authUser.Id,
		IsOtpValid: true,
	})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	output.UserId = verify.UserId
	output.IsValid = verify.IsValid
	output.PhoneNo = authUser.PhoneNo
	output.Note = otpUser.Note
	output.IsNewUser = authUser.PinCode == nil
	return
}
