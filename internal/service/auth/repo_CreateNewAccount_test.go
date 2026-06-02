package auth_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/dbtx"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
)

func TestAuthRepository_CreateNewAccount(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.CreateNewAccountInput
	}

	type expected = testutil.Result[auth.CreateNewAccountOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(tt *testing.T, mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "test failed to create new account",
			args: args{
				ctx:   context.Background(),
				input: auth.CreateNewAccountInput{},
			},
			expected: expected{
				Err:   auth.ErrCreateNewAccount,
				Value: auth.CreateNewAccountOutput{},
			},
			mock: func(tt *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrCreateNewAccount)
			},
		},
		{
			name: "test success to create new user",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewAccountInput{
					Id:      "358cbaad-316e-4539-9949-2636cdbd7e89",
					PhoneNo: "081234567890",
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.CreateNewAccountOutput{
					Id: "358cbaad-316e-4539-9949-2636cdbd7e89",
				},
			},
			mock: func(tt *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error init mock: %v", err)
			}
			if tt.mock != nil {
				tt.mock(t, sqlMock)
			}
			dbtx := dbtx.DbTx{Db: db}
			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: &dbtx,
				Db:            db,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})
			res, err := repo.CreateNewAccount(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
