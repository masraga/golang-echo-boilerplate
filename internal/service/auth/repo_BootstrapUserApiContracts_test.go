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

func TestAuthRepository_BootstrapUserApiContracts(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.BootstrapUserApiContractsInput
	}

	type expected = testutil.Result[auth.BootstrapUserApiContractsOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should return bootstrap error",
			args: args{
				ctx:   context.Background(),
				input: auth.BootstrapUserApiContractsInput{UserId: "358cbaad-316e-4539-9949-2636cdbd7e89"},
			},
			expected: expected{
				Err:   auth.ErrBootstrapUserApiContract,
				Value: auth.BootstrapUserApiContractsOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO public.auth_user_api_contract`).
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(auth.ErrBootstrapUserApiContract)
			},
		},
		{
			name: "should return inserted count",
			args: args{
				ctx:   context.Background(),
				input: auth.BootstrapUserApiContractsInput{UserId: "358cbaad-316e-4539-9949-2636cdbd7e89"},
			},
			expected: expected{
				Value: auth.BootstrapUserApiContractsOutput{InsertedCount: 3},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO public.auth_user_api_contract`).
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 3))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error initialize mock: %v", err)
			}
			defer dbMock.Close()

			if tt.mock != nil {
				tt.mock(sqlMock)
			}

			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: &dbtx.DbTx{Db: dbMock},
				Db:            dbMock,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})

			got, err := repo.BootstrapUserApiContracts(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
