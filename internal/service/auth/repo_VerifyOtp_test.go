package auth_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/leporo/sqlf"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/dbtx"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
)

func TestAuthRepository_VerifyOtp(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.VerifyOtpInput
	}

	type expected = testutil.Result[auth.VerifyOtpOutput]

	type fields struct {
		AuthRepoWriter *auth.MockAuthRepositoryWriterInterface
		DbTx           *dbtx.MockDbTxInterface
	}

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(sqlMock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "failed verify OTP",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyOtpInput{
					UserId:  "081234567890",
					OtpCode: "123456",
				},
			},
			expected: expected{
				Err:   auth.ErrVerifyOtp,
				Value: auth.VerifyOtpOutput{},
			},
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(``).WillReturnError(auth.ErrVerifyOtp)
			},
		},
		{
			name: "success verify OTP",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyOtpInput{
					UserId:  "081234567890",
					OtpCode: "123456",
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.VerifyOtpOutput{
					IsValid: true,
					UserId:  "081234567890",
				},
			},
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(``).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock, sqlMock, _ := sqlmock.New()
			tt.mock(sqlMock)

			dbTx := &dbtx.DbTx{Db: dbMock}
			defer dbMock.Close()

			authRepo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: dbTx,
				Sql:           sqlf.PostgreSQL,
				Db:            dbMock,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})

			got, err := authRepo.VerifyOtp(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
