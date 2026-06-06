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

func TestAuthRepository_VerifyUserAccount(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.VerifyUserAccountInput
	}

	type expected = testutil.Result[auth.VerifyUserAccountOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(sqlMock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "failed to verify user account",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyUserAccountInput{
					PhoneNo: "081234567890",
				},
			},
			expected: expected{
				Err:   auth.ErrVerifyUserAccount,
				Value: auth.VerifyUserAccountOutput{},
			},
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(``).WillReturnError(auth.ErrVerifyUserAccount)
			},
		},
		{
			name: "success to verify user account",
			args: args{
				ctx: context.Background(),
				input: auth.VerifyUserAccountInput{
					PhoneNo: "081234567890",
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.VerifyUserAccountOutput{
					IsVerified: true,
					PhoneNo:    "081234567890",
				},
			},
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(``).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock, _ := sqlmock.New()
			tt.mock(sqlMock)
			dbTx := &dbtx.DbTx{Db: db}
			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: dbTx,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				Sql:           sqlf.PostgreSQL,
			})

			got, err := repo.VerifyUserAccount(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
