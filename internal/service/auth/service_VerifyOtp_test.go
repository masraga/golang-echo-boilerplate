package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
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
				userId := faker.UUIDHyphenated()
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.
					EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      userId,
						PhoneNo: expectedPhoneNumber,
					}, nil)
				authRepoReader.
					EXPECT().
					FindOTP(gomock.Any(), gomock.Any()).
					Return(auth.FindOTPOutput{}, auth.ErrFindOTPNotFound)
				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     userId,
						IsOtpValid: false,
					}).
					Return(auth.UpdateOtpValidityOutput{}, nil)
				tt.fields.AuthRepoReader = authRepoReader
				tt.fields.AuthRepoWriter = authRepoWriter
			},
		},
		{
			name: "should fail when clearing otp validity fails",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyOtpInput{
					OtpCode: expectedOtp,
					PhoneNo: expectedPhoneNumber,
				},
			},
			expected: expected{Err: auth.ErrUpdateOtpValidity},
			mock: func(tt *test, ctrl *gomock.Controller) {
				userId := faker.UUIDHyphenated()
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      userId,
						PhoneNo: expectedPhoneNumber,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     userId,
						IsOtpValid: false,
					}).
					Return(auth.UpdateOtpValidityOutput{}, auth.ErrUpdateOtpValidity)

				tt.fields.AuthRepoReader = authRepoReader
				tt.fields.AuthRepoWriter = authRepoWriter
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
				userId := faker.UUIDHyphenated()
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.
					EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      userId,
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
				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     userId,
						IsOtpValid: false,
					}).
					Return(auth.UpdateOtpValidityOutput{}, nil)
				tt.fields.AuthRepoReader = authRepoReader
				tt.fields.AuthRepoWriter = authRepoWriter
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
				userId := faker.UUIDHyphenated()
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.
					EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      userId,
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
				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     userId,
						IsOtpValid: false,
					}).
					Return(auth.UpdateOtpValidityOutput{}, nil)
				tt.fields.AuthRepoReader = authRepoReader
				tt.fields.AuthRepoWriter = authRepoWriter
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
					IsValid:   true,
					UserId:    expectedPhoneNumber,
					PhoneNo:   expectedPhoneNumber,
					Note:      expectedNote,
					IsNewUser: true,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				userId := faker.UUIDHyphenated()
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.
					EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      userId,
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
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     userId,
						IsOtpValid: false,
					}).
					Return(auth.UpdateOtpValidityOutput{}, nil)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.
					EXPECT().
					VerifyOtp(gomock.Any(), gomock.Any()).
					Return(auth.VerifyOtpOutput{
						IsValid: true,
						UserId:  expectedPhoneNumber,
					}, nil)
				authRepoWriter.EXPECT().
					VerifyUserAccount(gomock.Any(), auth.VerifyUserAccountInput{
						PhoneNo: expectedPhoneNumber,
					}).
					Return(auth.VerifyUserAccountOutput{IsVerified: true}, nil)
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     userId,
						IsOtpValid: true,
					}).
					Return(auth.UpdateOtpValidityOutput{}, nil)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), nil).
					Return(nil)
				tt.fields.AuthRepoReader = authRepoReader
				tt.fields.AuthRepoWriter = authRepoWriter
			},
		},
		{
			name: "should keep otp validity false when otp persistence fails",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyOtpInput{
					OtpCode: expectedOtp,
					PhoneNo: expectedPhoneNumber,
				},
			},
			expected: expected{Err: auth.ErrVerifyOtp},
			mock: func(tt *test, ctrl *gomock.Controller) {
				userId := faker.UUIDHyphenated()
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      userId,
						PhoneNo: expectedPhoneNumber,
					}, nil)
				authRepoReader.EXPECT().
					FindOTP(gomock.Any(), gomock.Any()).
					Return(auth.FindOTPOutput{
						OtpCode:       expectedOtp,
						ExpiredAtUtc0: expectedValidExpiredAt,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     userId,
						IsOtpValid: false,
					}).
					Return(auth.UpdateOtpValidityOutput{}, nil)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					VerifyOtp(gomock.Any(), gomock.Any()).
					Return(auth.VerifyOtpOutput{}, auth.ErrVerifyOtp)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), auth.ErrVerifyOtp).
					Return(auth.ErrVerifyOtp)

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
