package auth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	"github.com/masraga/kerp-api/internal/util/pointer"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthService_AuthValidatePin(t *testing.T) {
	var (
		expectedUserId  = faker.UUIDHyphenated()
		expectedPhoneNo = "081234567890"
		expectedPinCode = "123456"
	)

	type args struct {
		ctx   context.Context
		input auth.AuthValidatePinInput
	}

	type expected = testutil.Result[auth.AuthValidatePinOutput]

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
				input: auth.AuthValidatePinInput{
					PhoneNo: expectedPhoneNo,
					PinCode: expectedPinCode,
				},
			},
			expected: expected{
				Err:   auth.ErrAuthNotFound,
				Value: auth.AuthValidatePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{}, auth.ErrAuthNotFound)

				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "should fail when otp is not validated",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo: expectedPhoneNo,
					PinCode: expectedPinCode,
				},
			},
			expected: expected{
				Err: auth.ErrOtpValidationRequired,
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						PinCode:    pointer.String(expectedPinCode),
						IsOtpValid: false,
					}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "should failed when new pin has no retype pin",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo: expectedPhoneNo,
					PinCode: expectedPinCode,
				},
			},
			expected: expected{
				Err:   auth.ErrValidateRetypePin,
				Value: auth.AuthValidatePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						IsOtpValid: true,
					}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "should failed when new pin and retype pin not match",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo:       expectedPhoneNo,
					PinCode:       expectedPinCode,
					RetypePinCode: pointer.String("654321"),
				},
			},
			expected: expected{
				Err:   auth.ErrPinCodeNotMatch,
				Value: auth.AuthValidatePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						IsOtpValid: true,
					}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "should failed when new pin length is invalid",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo:       expectedPhoneNo,
					PinCode:       "12345",
					RetypePinCode: pointer.String("12345"),
				},
			},
			expected: expected{
				Err:   auth.ErrPinIsTooLongOrShort,
				Value: auth.AuthValidatePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						IsOtpValid: true,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), auth.ErrPinIsTooLongOrShort).
					Return(auth.ErrPinIsTooLongOrShort)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should failed when begin db transaction failed",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo: expectedPhoneNo,
					PinCode: expectedPinCode,
				},
			},
			expected: expected{
				Err:   auth.ErrBeginDbTx,
				Value: auth.AuthValidatePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						PinCode:    pointer.String(expectedPinCode),
						IsOtpValid: true,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), auth.ErrBeginDbTx)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should failed when create new pin failed",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo:       expectedPhoneNo,
					PinCode:       expectedPinCode,
					RetypePinCode: pointer.String(expectedPinCode),
				},
			},
			expected: expected{
				Err:   auth.ErrCreateNewPin,
				Value: auth.AuthValidatePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						IsOtpValid: true,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					CreateNewPin(gomock.Any(), auth.CreateNewPinInput{
						UserId:  expectedUserId,
						PinCode: expectedPinCode,
					}).
					Return(auth.CreateNewPinOutput{}, auth.ErrCreateNewPin)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), auth.ErrCreateNewPin).
					Return(auth.ErrCreateNewPin)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should failed when store access token failed",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo: expectedPhoneNo,
					PinCode: expectedPinCode,
				},
			},
			expected: expected{
				Err:   auth.ErrStoreAccessToken,
				Value: auth.AuthValidatePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						PinCode:    pointer.String(expectedPinCode),
						IsOtpValid: true,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					StoreAccessToken(gomock.Any(), gomock.Any()).
					Return(auth.StoreAccessTokenOutput{}, auth.ErrStoreAccessToken)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), auth.ErrStoreAccessToken).
					Return(auth.ErrStoreAccessToken)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should failed when existing pin not match",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo: expectedPhoneNo,
					PinCode: "654321",
				},
			},
			expected: expected{
				Err:   auth.ErrPinCodeNotMatch,
				Value: auth.AuthValidatePinOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						PinCode:    pointer.String(expectedPinCode),
						IsOtpValid: true,
					}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "should success when creating new pin",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo:       expectedPhoneNo,
					PinCode:       expectedPinCode,
					RetypePinCode: pointer.String(expectedPinCode),
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.AuthValidatePinOutput{
					IsValid: true,
					UserId:  expectedUserId,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						IsOtpValid: true,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					CreateNewPin(gomock.Any(), auth.CreateNewPinInput{
						UserId:  expectedUserId,
						PinCode: expectedPinCode,
					}).
					Return(auth.CreateNewPinOutput{
						UserId:  expectedUserId,
						PinCode: expectedPinCode,
					}, nil)
				authRepoWriter.EXPECT().
					StoreAccessToken(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, input auth.StoreAccessTokenInput) (auth.StoreAccessTokenOutput, error) {
						require.Equal(t, expectedUserId, input.UserId)
						require.NotEmpty(t, input.Token)
						require.Greater(t, input.ExpiredAtUtc0, int64(0))
						return auth.StoreAccessTokenOutput{
							Token:         input.Token,
							UserId:        input.UserId,
							ExpiredAtUtc0: input.ExpiredAtUtc0,
							IsActive:      true,
						}, nil
					})
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     expectedUserId,
						IsOtpValid: false,
					}).
					Return(auth.UpdateOtpValidityOutput{}, nil)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), nil).
					Return(nil)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should success when validating existing pin",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo: expectedPhoneNo,
					PinCode: expectedPinCode,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.AuthValidatePinOutput{
					IsValid: true,
					UserId:  expectedUserId,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						PinCode:    pointer.String(expectedPinCode),
						IsOtpValid: true,
					}, nil)

				tt.fields.AuthRepositoryReader = authRepoReader
				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					StoreAccessToken(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, input auth.StoreAccessTokenInput) (auth.StoreAccessTokenOutput, error) {
						require.Equal(t, expectedUserId, input.UserId)
						require.NotEmpty(t, input.Token)
						require.Greater(t, input.ExpiredAtUtc0, int64(0))
						return auth.StoreAccessTokenOutput{
							Token:         input.Token,
							UserId:        input.UserId,
							ExpiredAtUtc0: input.ExpiredAtUtc0,
							IsActive:      true,
						}, nil
					})
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     expectedUserId,
						IsOtpValid: false,
					}).
					Return(auth.UpdateOtpValidityOutput{}, nil)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), nil).
					Return(nil)

				tt.fields.AuthRepositoryWriter = authRepoWriter
			},
		},
		{
			name: "should rollback when consuming otp validity fails",
			args: args{
				ctx: context.Background(),
				input: auth.AuthValidatePinInput{
					PhoneNo: expectedPhoneNo,
					PinCode: expectedPinCode,
				},
			},
			expected: expected{Err: auth.ErrUpdateOtpValidity},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), auth.FindAuthInput{PhoneNo: expectedPhoneNo}).
					Return(auth.FindAuthOutput{
						Id:         expectedUserId,
						PhoneNo:    expectedPhoneNo,
						PinCode:    pointer.String(expectedPinCode),
						IsOtpValid: true,
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(gomock.Any(), nil).
					Return(context.Background(), nil)
				authRepoWriter.EXPECT().
					StoreAccessToken(gomock.Any(), gomock.Any()).
					Return(auth.StoreAccessTokenOutput{}, nil)
				authRepoWriter.EXPECT().
					UpdateOtpValidity(gomock.Any(), auth.UpdateOtpValidityInput{
						UserId:     expectedUserId,
						IsOtpValid: false,
					}).
					Return(auth.UpdateOtpValidityOutput{}, auth.ErrUpdateOtpValidity)
				authRepoWriter.EXPECT().
					CommitOrRollback(gomock.Any(), auth.ErrUpdateOtpValidity).
					Return(auth.ErrUpdateOtpValidity)

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

			svc := auth.NewAuthService(auth.AuthServiceOpts{
				JwtSecret:            "test-secret",
				JwtExpiration:        15,
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader: tt.fields.AuthRepositoryReader,
				AuthRepositoryWriter: tt.fields.AuthRepositoryWriter,
			})

			got, err := svc.AuthValidatePin(tt.args.ctx, tt.args.input)
			if tt.expected.Err != nil {
				testutil.RequireResult(t, err, tt.expected, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected.Value.IsValid, got.IsValid)
			require.Equal(t, tt.expected.Value.UserId, got.UserId)
			require.NotEmpty(t, got.Token)
			require.Greater(t, got.ExpiredAtUtc0, int64(0))
		})
	}
}
