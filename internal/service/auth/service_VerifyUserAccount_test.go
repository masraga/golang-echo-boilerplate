package auth_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/dbtx"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	gomock "go.uber.org/mock/gomock"
)

func TestAuthService_VerifyUserAccount(t *testing.T) {
	var (
		expectedUserId    string = faker.UUIDHyphenated()
		expectedPhoneNo   string = "081234567890"
		expectedIsNewUser bool   = true
	)

	type args struct {
		ctx   context.Context
		input auth.VerifyUserAccountInput
	}

	type fields struct {
		DbTxInterface  *dbtx.DbTxInterface
		AuthRepoReader *auth.MockAuthRepositoryReaderInterface
		AuthRepoWriter *auth.MockAuthRepositoryWriterInterface
	}

	type expected = testutil.Result[auth.VerifyUserAccountOutput]

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "failed to find auth user",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyUserAccountInput{
					PhoneNo: expectedPhoneNo,
				},
			},
			expected: expected{
				Err:   auth.ErrAuthNotFound,
				Value: auth.VerifyUserAccountOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{}, auth.ErrAuthNotFound)

				tt.fields.AuthRepoReader = authRepoReader
			},
		},
		{
			name: "failed to verify account",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyUserAccountInput{
					PhoneNo: expectedPhoneNo,
				},
			},
			expected: expected{
				Err:   auth.ErrVerifyUserAccount,
				Value: auth.VerifyUserAccountOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{PhoneNo: expectedPhoneNo, Id: expectedUserId}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					VerifyUserAccount(gomock.Any(), gomock.Any()).
					Return(auth.VerifyUserAccountOutput{}, auth.ErrVerifyUserAccount)

				tt.fields.AuthRepoReader = authRepoReader
				tt.fields.AuthRepoWriter = authRepoWriter
			},
		},
		{
			name: "success to verify new user account",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyUserAccountInput{
					PhoneNo: expectedPhoneNo,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.VerifyUserAccountOutput{
					UserId:    expectedUserId,
					PhoneNo:   expectedPhoneNo,
					IsNewUser: expectedIsNewUser,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					FindAuth(gomock.Any(), gomock.Any()).
					Return(auth.FindAuthOutput{PhoneNo: expectedPhoneNo, Id: expectedUserId}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					VerifyUserAccount(gomock.Any(), gomock.Any()).
					Return(auth.VerifyUserAccountOutput{
						UserId:    expectedUserId,
						PhoneNo:   expectedPhoneNo,
						IsNewUser: expectedIsNewUser,
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

			tt.mock(&tt, ctrl)

			s := auth.NewAuthService(auth.AuthServiceOpts{
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader: tt.fields.AuthRepoReader,
				AuthRepositoryWriter: tt.fields.AuthRepoWriter,
			})

			got, err := s.VerifyUserAccount(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
