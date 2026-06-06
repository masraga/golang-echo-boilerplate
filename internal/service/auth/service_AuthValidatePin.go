package auth

import (
	"context"
	"errors"
	"time"

	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
)

func (s *AuthService) AuthValidatePin(ctx context.Context, input AuthValidatePinInput) (output AuthValidatePinOutput, err error) {
	authUser, err := s.AuthRepositoryReader.FindAuth(ctx, FindAuthInput{
		PhoneNo: input.PhoneNo,
	})
	if err != nil {
		return
	}
	if !authUser.IsOtpValid {
		err = s.Err.Wrap(ErrOtpValidationRequired)
		return
	}
	if authUser.PinCode == nil && input.RetypePinCode == nil {
		err = s.Err.Wrap(ErrValidateRetypePin)
		return
	}
	if authUser.PinCode == nil && (input.PinCode != pointer.SafeString(input.RetypePinCode)) {
		err = s.Err.Wrap(ErrPinCodeNotMatch)
		return
	}
	if authUser.PinCode != nil && pointer.SafeString(authUser.PinCode) != input.PinCode {
		err = s.Err.Wrap(ErrPinCodeNotMatch)
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

	//create new pin is not set
	input.UserId = authUser.Id
	if authUser.PinCode == nil {
		output, err = s.createNewPin(ctx, input)
		authUser.PinCode = pointer.String(input.PinCode)
		if err != nil {
			return
		}
	}

	// create new token if pin code is valid
	expiredAtUtc0 := time.Now().Add(time.Duration(JwtTokenExpiredDuration) * time.Minute).UnixMilli()
	token, err := s.CreateToken(ctx, UserTokenClaimInput{
		TokenType:     TokenTypeJwt,
		UserId:        authUser.Id,
		ExpiredAtUtc0: expiredAtUtc0,
		IssuerAtUtc0:  time.Now().UnixMilli(),
		UserName:      authUser.PhoneNo,
	})
	if err != nil {
		err = s.Err.Wrap(err)
		return
	}

	_, err = s.AuthRepositoryWriter.StoreAccessToken(ctx, StoreAccessTokenInput{
		Token:         token.Token,
		UserId:        authUser.Id,
		ExpiredAtUtc0: expiredAtUtc0,
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

	output.UserId = authUser.Id
	output.IsValid = true
	output.Token = token.Token
	output.ExpiredAtUtc0 = expiredAtUtc0

	return
}

func (s *AuthService) createNewPin(ctx context.Context, input AuthValidatePinInput) (output AuthValidatePinOutput, err error) {
	if len(input.PinCode) < MinPinLen || len(input.PinCode) > MaxPinLen {
		err = s.Err.Wrap(ErrPinIsTooLongOrShort)
		return
	}

	authPin, err := s.AuthRepositoryWriter.CreateNewPin(ctx, CreateNewPinInput{
		PinCode: input.PinCode,
		UserId:  input.UserId,
	})
	if err != nil {
		return
	}

	output.IsValid = authPin.PinCode == input.PinCode
	return
}
