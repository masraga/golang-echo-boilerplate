package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
)

func (s *AuthService) CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error) {
	authUser, err := s.AuthRepositoryReader.FindAuth(ctx, FindAuthInput{
		PhoneNo: input.PhoneNo,
	})
	if err != nil && !errors.Is(err, ErrAuthNotFound) {
		err = s.Err.Wrap(err)
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

	userId := authUser.Id
	if userId == "" {
		input.Id = uuid.NewString()

		_, err = s.AuthRepositoryWriter.CreateNewAccount(ctx, input)
		if err != nil {
			err = s.Err.Wrap(err)
			return
		}

		userId = input.Id
	} else {
		_, err = s.AuthRepositoryWriter.UpdateFirebaseId(ctx, UpdateFirebaseIdInput{
			UserId:     userId,
			FirebaseId: *input.FirebaseId,
		})
		if err != nil {
			err = s.Err.Wrap(err)
			return
		}

		_, err = s.AuthRepositoryWriter.UpdateOtpValidity(ctx, UpdateOtpValidityInput{
			UserId:     userId,
			IsOtpValid: false,
		})
		if err != nil {
			err = s.Err.Wrap(err)
			return
		}
	}

	otpSvc, err := s.createRegistrationOTP(ctx, CreateOTPInput{
		UserId:  userId,
		PhoneNo: input.PhoneNo,
	})
	if err != nil {
		return
	}

	output.Id = userId
	output.OtpCode = otpSvc.OtpCode
	return
}

func (s *AuthService) createRegistrationOTP(ctx context.Context, input CreateOTPInput) (output CreateOTPOutput, err error) {
	input.Note = pointer.String("OTP for account register")
	return s.CreateOTP(ctx, input)
}
