package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	"github.com/masraga/kerp-api/internal/util/pointer"
	"go.uber.org/mock/gomock"
)

func TestAuthService_VerifyOtp(t *testing.T) {
	var (
		expectedPhoneNumber    string  = "081234567890"
		expectedOtp            string  = "454546"
		expectedNote           *string = pointer.String(faker.Word())
		expectedExpiredAt      int64   = time.Now().Add(-15 * time.Minute).UnixMilli()
		expectedValidExpiredAt int64   = time.Now().Add(15 * time.Minute).UnixMilli()
	)
	type args struct {
		ctx   context.Context
		input auth.VerifyOtpInput
	}

	type fields struct {
		AuthRepoReader *auth.MockAuthRepositoryReaderInterface
		AuthRepoWriter *auth.MockAuthRepositoryWriterInterface
	}

	type expected = testutil.Result[auth.VerifyOtpOutput]

	type test struct {
		name     string
		args     args
		expected expected
		fields   fields
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should fail when find auth",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyOtpInput{
					OtpCode: expectedOtp,
					UserId:  expectedPhoneNumber,
				},
			},
			expected: expected{
				Err:   auth.ErrAuthNotFound,
				Value: auth.VerifyOtpOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.
					EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{}, auth.ErrAuthNotFound)
				tt.fields.AuthRepoReader = authRepoReader
			},
		},
		{
			name: "should fail when find otp",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyOtpInput{
					OtpCode: expectedOtp,
					UserId:  expectedPhoneNumber,
				},
			},
			expected: expected{
				Err:   auth.ErrFindOTPNotFound,
				Value: auth.VerifyOtpOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.
					EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      faker.UUIDHyphenated(),
						PhoneNo: expectedPhoneNumber,
					}, nil)
				authRepoReader.
					EXPECT().
					FindOTP(gomock.Any(), gomock.Any()).
					Return(auth.FindOTPOutput{}, auth.ErrFindOTPNotFound)
				tt.fields.AuthRepoReader = authRepoReader
			},
		},
		{
			name: "should fail when otp is expired",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyOtpInput{
					OtpCode: expectedOtp,
					UserId:  expectedPhoneNumber,
				},
			},
			expected: expected{
				Err:   auth.ErrOtpIsExpired,
				Value: auth.VerifyOtpOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.
					EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      faker.UUIDHyphenated(),
						PhoneNo: expectedPhoneNumber,
					}, nil)
				authRepoReader.
					EXPECT().
					FindOTP(gomock.Any(), gomock.Any()).
					Return(auth.FindOTPOutput{
						OtpCode:       expectedOtp,
						Note:          expectedNote,
						ExpiredAtUtc0: expectedExpiredAt,
					}, nil)
				tt.fields.AuthRepoReader = authRepoReader
			},
		},
		{
			name: "should fail when OTP already verified",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyOtpInput{
					OtpCode: expectedOtp,
					UserId:  expectedPhoneNumber,
				},
			},
			expected: expected{
				Err:   auth.ErrOtpAleadyVerified,
				Value: auth.VerifyOtpOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.
					EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      faker.UUIDHyphenated(),
						PhoneNo: expectedPhoneNumber,
					}, nil)
				authRepoReader.
					EXPECT().
					FindOTP(gomock.Any(), gomock.Any()).
					Return(auth.FindOTPOutput{
						OtpCode:       expectedOtp,
						Note:          expectedNote,
						ExpiredAtUtc0: expectedValidExpiredAt,
						IsVerified:    true,
					}, nil)
				tt.fields.AuthRepoReader = authRepoReader
			},
		},
		{
			name: "should success verify otp",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyOtpInput{
					OtpCode: expectedOtp,
					UserId:  expectedPhoneNumber,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.VerifyOtpOutput{
					IsValid: true,
					UserId:  expectedPhoneNumber,
					Note:    expectedNote,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.
					EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      faker.UUIDHyphenated(),
						PhoneNo: expectedPhoneNumber,
					}, nil)
				authRepoReader.
					EXPECT().
					FindOTP(gomock.Any(), gomock.Any()).
					Return(auth.FindOTPOutput{
						OtpCode:       expectedOtp,
						Note:          expectedNote,
						ExpiredAtUtc0: expectedValidExpiredAt,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.
					EXPECT().
					VerifyOtp(gomock.Any(), gomock.Any()).
					Return(auth.VerifyOtpOutput{
						IsValid: true,
						UserId:  expectedPhoneNumber,
					}, nil)
				tt.fields.AuthRepoReader = authRepoReader
				tt.fields.AuthRepoWriter = authRepoWriter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.mock != nil {
				tt.mock(&tt, ctrl)
			}

			authService := auth.NewAuthService(auth.AuthServiceOpts{
				JwtSecret:            auth.JwtSecretType(faker.Word()),
				JwtExpiration:        auth.JwtExpirationType(60),
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader: tt.fields.AuthRepoReader,
				AuthRepositoryWriter: tt.fields.AuthRepoWriter,
			})

			got, err := authService.VerifyOtp(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
