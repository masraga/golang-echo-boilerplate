package auth_test

import (
	"context"
	"testing"

	"github.com/masraga/kerp-api/internal/ctxerr"
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

	type expected struct {
		Err    error
		Value  auth.CreateNewAccountOutput
		Assert func(t *testing.T, actual auth.CreateNewAccountOutput)
	}

	type fields struct {
		AuthRepositoryReader *auth.MockAuthRepositoryReaderInterface
		AuthRepositoryWriter *auth.MockAuthRepositoryWriterInterface
	}

	type test struct {
		name     string
		args     args
		expected expected
		fields   fields
		mock     func(t *testing.T, tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should create otp and update firebase id for existing user",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewAccountInput{
					PhoneNo:    "081234567890",
					FirebaseId: stringPointer("fcm-registration-token"),
				},
			},
			expected: expected{
				Value: auth.CreateNewAccountOutput{
					Id: "358cbaad-316e-4539-9949-2636cdbd7e89",
				},
				Assert: func(t *testing.T, actual auth.CreateNewAccountOutput) {
					require.Equal(t, "358cbaad-316e-4539-9949-2636cdbd7e89", actual.Id)
					require.Regexp(t, `^[0-9]{6}$`, actual.OtpCode)
				},
			},
			mock: func(t *testing.T, tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: tt.args.input.PhoneNo}).
					Return(auth.FindAuthOutput{
						Id:      tt.expected.Value.Id,
						PhoneNo: tt.args.input.PhoneNo,
					}, nil)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{
						PhoneNo: tt.args.input.PhoneNo,
						UserId:  tt.expected.Value.Id,
					}).
					Return(auth.FindAuthOutput{
						Id:      tt.expected.Value.Id,
						PhoneNo: tt.args.input.PhoneNo,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					UpdateFirebaseId(gomock.Any(), auth.UpdateFirebaseIdInput{
						UserId:     tt.expected.Value.Id,
						FirebaseId: *tt.args.input.FirebaseId,
					}).
					Return(auth.UpdateFirebaseIdOutput{}, nil)
				authRepoWriter.EXPECT().
					DeleteAllUserOTP(gomock.Any(), auth.DeleteAllUserOTPInput{UserId: tt.expected.Value.Id}).
					Return(auth.DeleteAllUserOTPOutput{IsSuccess: true}, nil)
				authRepoWriter.EXPECT().
					CreateOTP(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, input auth.CreateOTPInput) (auth.CreateOTPOutput, error) {
						require.Equal(t, tt.expected.Value.Id, input.UserId)
						require.NotNil(t, input.Note)
						require.Equal(t, "OTP for account register", *input.Note)
						return auth.CreateOTPOutput{OtpCode: input.OtpCode}, nil
					})
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), nil).
					Return(nil)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should update firebase id for existing user",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewAccountInput{
					PhoneNo:    "081234567890",
					FirebaseId: stringPointer("new-fcm-registration-token"),
				},
			},
			expected: expected{
				Value: auth.CreateNewAccountOutput{
					Id: "358cbaad-316e-4539-9949-2636cdbd7e89",
				},
				Assert: func(t *testing.T, actual auth.CreateNewAccountOutput) {
					require.Equal(t, "358cbaad-316e-4539-9949-2636cdbd7e89", actual.Id)
					require.Regexp(t, `^[0-9]{6}$`, actual.OtpCode)
				},
			},
			mock: func(t *testing.T, tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: tt.args.input.PhoneNo}).
					Return(auth.FindAuthOutput{Id: tt.expected.Value.Id}, nil)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{
						PhoneNo: tt.args.input.PhoneNo,
						UserId:  tt.expected.Value.Id,
					}).
					Return(auth.FindAuthOutput{Id: tt.expected.Value.Id}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					UpdateFirebaseId(gomock.Any(), auth.UpdateFirebaseIdInput{
						UserId:     tt.expected.Value.Id,
						FirebaseId: *tt.args.input.FirebaseId,
					}).
					Return(auth.UpdateFirebaseIdOutput{}, nil)
				authRepoWriter.EXPECT().
					DeleteAllUserOTP(gomock.Any(), auth.DeleteAllUserOTPInput{UserId: tt.expected.Value.Id}).
					Return(auth.DeleteAllUserOTPOutput{IsSuccess: true}, nil)
				authRepoWriter.EXPECT().
					CreateOTP(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, input auth.CreateOTPInput) (auth.CreateOTPOutput, error) {
						return auth.CreateOTPOutput{OtpCode: input.OtpCode}, nil
					})
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), nil).
					Return(nil)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should rollback when firebase id update fails",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewAccountInput{
					PhoneNo:    "081234567890",
					FirebaseId: stringPointer("new-fcm-registration-token"),
				},
			},
			expected: expected{
				Err: auth.ErrUpdateFirebaseId,
			},
			mock: func(t *testing.T, tt *test, ctrl *gomock.Controller) {
				userId := "358cbaad-316e-4539-9949-2636cdbd7e89"
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: tt.args.input.PhoneNo}).
					Return(auth.FindAuthOutput{Id: userId}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					UpdateFirebaseId(gomock.Any(), auth.UpdateFirebaseIdInput{
						UserId:     userId,
						FirebaseId: *tt.args.input.FirebaseId,
					}).
					Return(auth.UpdateFirebaseIdOutput{}, auth.ErrUpdateFirebaseId)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), auth.ErrUpdateFirebaseId).
					Return(auth.ErrUpdateFirebaseId)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should failed begin dbtx",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewAccountInput{
					PhoneNo:    "081234567890",
					FirebaseId: stringPointer("fcm-registration-token"),
				},
			},
			expected: expected{
				Err:   auth.ErrBeginDbTx,
				Value: auth.CreateNewAccountOutput{},
			},
			mock: func(t *testing.T, tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: tt.args.input.PhoneNo}).
					Return(auth.FindAuthOutput{}, auth.ErrAuthNotFound)

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
					PhoneNo:    "081234567890",
					FirebaseId: stringPointer("fcm-registration-token"),
				},
			},
			expected: expected{
				Err:   auth.ErrCreateNewOTP,
				Value: auth.CreateNewAccountOutput{},
			},
			mock: func(t *testing.T, tt *test, ctrl *gomock.Controller) {
				var createdAccount auth.CreateNewAccountInput

				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: tt.args.input.PhoneNo}).
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
					Return(auth.DeleteAllUserOTPOutput{IsSuccess: true}, nil)
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
		{
			name: "should return account id and persisted otp when successful",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewAccountInput{
					PhoneNo:    "081234567890",
					FirebaseId: stringPointer("fcm-registration-token"),
				},
			},
			expected: expected{
				Assert: func(t *testing.T, actual auth.CreateNewAccountOutput) {
					require.NotEmpty(t, actual.Id)
					require.Regexp(t, `^[0-9]{6}$`, actual.OtpCode)
				},
			},
			mock: func(t *testing.T, tt *test, ctrl *gomock.Controller) {
				var createdAccount auth.CreateNewAccountInput
				var createdOTP auth.CreateOTPInput

				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: tt.args.input.PhoneNo}).
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
						require.Equal(t, tt.args.input.FirebaseId, input.FirebaseId)
						return auth.CreateNewAccountOutput{Id: input.Id}, nil
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

				tt.expected.Assert = func(t *testing.T, actual auth.CreateNewAccountOutput) {
					require.Equal(t, createdAccount.Id, actual.Id)
					require.Equal(t, createdOTP.OtpCode, actual.OtpCode)
					require.Regexp(t, `^[0-9]{6}$`, actual.OtpCode)
				}

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
				tt.mock(t, &tt, ctrl)
			}

			authService := auth.NewAuthService(auth.AuthServiceOpts{
				JwtSecret:            "secret",
				JwtExpiration:        3600,
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader: tt.fields.AuthRepositoryReader,
				AuthRepositoryWriter: tt.fields.AuthRepositoryWriter,
			})

			res, err := authService.CreateNewAccount(tt.args.ctx, tt.args.input)
			if tt.expected.Assert != nil {
				require.NoError(t, err)
				tt.expected.Assert(t, res)
				return
			}

			testutil.RequireResult(t, err, testutil.Result[auth.CreateNewAccountOutput]{
				Err:   tt.expected.Err,
				Value: tt.expected.Value,
			}, res)
		})
	}
}

func stringPointer(value string) *string {
	return &value
}
