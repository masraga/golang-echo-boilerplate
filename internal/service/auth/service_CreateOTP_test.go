package auth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
	"go.uber.org/mock/gomock"
)

func TestAuthService_CreateOTP(t *testing.T) {
	var (
		expectedUserId  string  = faker.UUIDHyphenated()
		expectedOtpCode string  = "123456"
		expectedNote    *string = pointer.String(faker.Word())
	)

	type args struct {
		ctx   context.Context
		input auth.CreateOTPInput
	}

	type expected = testutil.Result[auth.CreateOTPOutput]

	type fields struct {
		AuthRepositoryReader *auth.MockAuthRepositoryReaderInterface
		AuthRepositoryWriter *auth.MockAuthRepositoryWriterInterface
	}

	type test struct {
		name     string
		args     args
		expected expected
		fields   fields
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should failed when auth user not found",
			args: args{
				ctx: context.Background(),
				input: auth.CreateOTPInput{
					UserId:  expectedUserId,
					Note:    expectedNote,
					OtpCode: expectedOtpCode,
				},
			},
			expected: expected{
				Err:   auth.ErrAuthNotFound,
				Value: auth.CreateOTPOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepo := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepo.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{}, auth.ErrAuthNotFound)
				tt.fields.AuthRepositoryReader = authRepo
			},
		},
		{
			name: "should failed when delete all user otp fails",
			args: args{
				ctx: context.Background(),
				input: auth.CreateOTPInput{
					UserId:  expectedUserId,
					Note:    expectedNote,
					OtpCode: expectedOtpCode,
				},
			},
			expected: expected{
				Err:   auth.ErrFailedDeleteUserOTP,
				Value: auth.CreateOTPOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id: expectedUserId,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					DeleteAllUserOTP(gomock.Any(), gomock.Any()).
					Return(auth.DeleteAllUserOTPOutput{}, auth.ErrFailedDeleteUserOTP)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should failed when create otp",
			args: args{
				ctx: context.Background(),
				input: auth.CreateOTPInput{
					UserId:  expectedUserId,
					Note:    expectedNote,
					OtpCode: expectedOtpCode,
				},
			},
			expected: expected{
				Err:   auth.ErrCreateNewOTP,
				Value: auth.CreateOTPOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id: expectedUserId,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					DeleteAllUserOTP(gomock.Any(), gomock.Any()).
					Return(auth.DeleteAllUserOTPOutput{IsSuccess: true}, nil)
				authRepoWriter.EXPECT().
					CreateOTP(gomock.Any(), gomock.Any()).
					Return(auth.CreateOTPOutput{}, auth.ErrCreateNewOTP)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should success creating OTP",
			args: args{
				ctx: context.Background(),
				input: auth.CreateOTPInput{
					UserId:  expectedUserId,
					Note:    expectedNote,
					OtpCode: expectedOtpCode,
				},
			},
			expected: expected{
				Err:   nil,
				Value: auth.CreateOTPOutput{OtpCode: expectedOtpCode},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id: expectedUserId,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					DeleteAllUserOTP(gomock.Any(), gomock.Any()).
					Return(auth.DeleteAllUserOTPOutput{IsSuccess: true}, nil)
				authRepoWriter.EXPECT().
					CreateOTP(gomock.Any(), gomock.Any()).
					Return(auth.CreateOTPOutput{OtpCode: expectedOtpCode}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.mock(&tt, ctrl)

			s := auth.NewAuthService(auth.AuthServiceOpts{
				JwtSecret:            "test-secret",
				JwtExpiration:        15,
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader: tt.fields.AuthRepositoryReader,
				AuthRepositoryWriter: tt.fields.AuthRepositoryWriter,
			})

			got, err := s.CreateOTP(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
