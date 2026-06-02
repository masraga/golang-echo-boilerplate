package auth

import (
	"context"

	"github.com/masraga/kerp-api/internal/dbtx"
)

type AuthServiceInterface interface {
	CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error)
	CreateToken(ctx context.Context, input UserTokenClaimInput) (output UserTokenClaimOutput, err error)
	CreateJWTToken(ctx context.Context, input CreateJWTTokenInput) (output CreateJWTTokenOutput, err error)
	CreateOTP(ctx context.Context, input CreateOTPInput) (output CreateOTPOutput, err error)
	VerifyOtp(ctx context.Context, input VerifyOtpInput) (output VerifyOtpOutput, err error)
	VerifyUserAccount(ctx context.Context, input VerifyUserAccountInput) (output VerifyUserAccountOutput, err error)
	AuthValidatePin(ctx context.Context, input AuthValidatePinInput) (output AuthValidatePinOutput, err error)
	ValidateJwtToken(ctx context.Context, input ValidateJwtTokenInput) (output ValidateJwtTokenOutput, err error)
}

type AuthRepositoryWriterInterface interface {
	dbtx.DbTxInterface
	CreateNewAccount(ctx context.Context, input CreateNewAccountInput) (output CreateNewAccountOutput, err error)
	DeleteAllUserOTP(ctx context.Context, input DeleteAllUserOTPInput) (output DeleteAllUserOTPOutput, err error)
	CreateOTP(ctx context.Context, input CreateOTPInput) (output CreateOTPOutput, err error)
	VerifyOtp(ctx context.Context, input VerifyOtpInput) (output VerifyOtpOutput, err error)
	VerifyUserAccount(ctx context.Context, input VerifyUserAccountInput) (output VerifyUserAccountOutput, err error)
	CreateNewPin(ctx context.Context, input CreateNewPinInput) (output CreateNewPinOutput, err error)
	StoreAccessToken(ctx context.Context, input StoreAccessTokenInput) (output StoreAccessTokenOutput, err error)
}

type AuthRepositoryReaderInterface interface {
	FindAuth(ctx context.Context, input FindAuthInput) (output FindAuthOutput, err error)
	FindOTP(ctx context.Context, input FindOTPInput) (output FindOTPOutput, err error)
	FindAccessToken(ctx context.Context, input FindAccessTokenInput) (output FindAccessTokenOutput, err error)
}
