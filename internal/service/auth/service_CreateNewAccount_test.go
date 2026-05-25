package auth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthService_CreateNewAccount(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.CreateNewAccountInput
	}

	type expected = testutil.Result[auth.CreateNewAccountOutput]

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
			name: "should failed when user already exist",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewAccountInput{
					PhoneNo: "081234567890",
				},
			},
			expected: expected{
				Err:   auth.ErrDuplicateUser,
				Value: auth.CreateNewAccountOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{
						Id:      faker.UUIDHyphenated(),
						PhoneNo: "081234567890",
					}, nil)
				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "should failed begin dbtx",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewAccountInput{
					PhoneNo: "081234567890",
				},
			},
			expected: expected{
				Err:   auth.ErrBeginDbTx,
				Value: auth.CreateNewAccountOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), auth.ErrBeginDbTx)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should rollback when create otp failed",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewAccountInput{
					PhoneNo: "081234567890",
				},
			},
			expected: expected{
				Err:   auth.ErrCreateNewOTP,
				Value: auth.CreateNewAccountOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{}, auth.ErrAuthNotFound)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{Id: faker.UUIDHyphenated()}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					CreateNewAccount(gomock.Any(), gomock.Any()).
					Return(auth.CreateNewAccountOutput{}, nil)
				authRepoWriter.EXPECT().
					DeleteAllUserOTP(gomock.Any(), gomock.Any()).
					Return(auth.DeleteAllUserOTPOutput{}, nil)
				authRepoWriter.EXPECT().
					CreateOTP(gomock.Any(), gomock.Any()).
					Return(auth.CreateOTPOutput{}, auth.ErrCreateNewOTP)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), auth.ErrCreateNewOTP).
					Return(auth.ErrCreateNewOTP)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
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
				JwtSecret:            "secret",
				JwtExpiration:        3600,
				AuthRepositoryReader: tt.fields.AuthRepositoryReader,
				AuthRepositoryWriter: tt.fields.AuthRepositoryWriter,
			})

			res, err := authService.CreateNewAccount(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}

	t.Run("should return account id and persisted otp when successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var createdAccount auth.CreateNewAccountInput
		var createdOTP auth.CreateOTPInput

		authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
		authRepoReader.EXPECT().
			FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: "081234567890"}).
			Return(auth.FindAuthOutput{}, auth.ErrAuthNotFound)
		authRepoReader.EXPECT().
			FindAuth(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, input auth.FindAuthInput) (auth.FindAuthOutput, error) {
				require.Equal(t, createdAccount.Id, input.UserId)
				return auth.FindAuthOutput{Id: input.UserId}, nil
			})

		authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
		authRepoWriter.EXPECT().
			Begin(gomock.Any(), nil).
			Return(context.Background(), nil)
		authRepoWriter.EXPECT().
			CreateNewAccount(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, input auth.CreateNewAccountInput) (auth.CreateNewAccountOutput, error) {
				createdAccount = input
				return auth.CreateNewAccountOutput{}, nil
			})
		authRepoWriter.EXPECT().
			DeleteAllUserOTP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, input auth.DeleteAllUserOTPInput) (auth.DeleteAllUserOTPOutput, error) {
				require.Equal(t, createdAccount.Id, input.UserId)
				return auth.DeleteAllUserOTPOutput{IsSuccess: true}, nil
			})
		authRepoWriter.EXPECT().
			CreateOTP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, input auth.CreateOTPInput) (auth.CreateOTPOutput, error) {
				createdOTP = input
				return auth.CreateOTPOutput{OtpCode: input.OtpCode}, nil
			})
		authRepoWriter.EXPECT().
			CommitOrRollback(gomock.Any(), nil).
			Return(nil)

		authService := auth.NewAuthService(auth.AuthServiceOpts{
			JwtSecret:            "secret",
			JwtExpiration:        3600,
			AuthRepositoryReader: authRepoReader,
			AuthRepositoryWriter: authRepoWriter,
		})

		res, err := authService.CreateNewAccount(context.Background(), auth.CreateNewAccountInput{
			PhoneNo: "081234567890",
		})

		require.NoError(t, err)
		require.NotEmpty(t, res.Id)
		require.Regexp(t, `^[0-9]{6}$`, res.OtpCode)
		require.Equal(t, createdAccount.Id, res.Id)
		require.Equal(t, createdOTP.OtpCode, res.OtpCode)
	})
}
