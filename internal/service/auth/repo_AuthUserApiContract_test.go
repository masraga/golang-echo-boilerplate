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

func TestAuthRepository_DeleteAuthUserApiContract(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.DeleteAuthUserApiContractInput
	}

	type expected = testutil.Result[auth.DeleteAuthUserApiContractOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should hard delete auth user api contract",
			args: args{
				ctx:   context.Background(),
				input: auth.DeleteAuthUserApiContractInput{Id: "1d7e2f23-b4c3-4ad7-8478-e770b68e6f11"},
			},
			expected: expected{
				Value: auth.DeleteAuthUserApiContractOutput{IsSuccess: true},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM public.auth_user_api_contract`).
					WithArgs("1d7e2f23-b4c3-4ad7-8478-e770b68e6f11").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "should return not found when no row deleted",
			args: args{
				ctx:   context.Background(),
				input: auth.DeleteAuthUserApiContractInput{Id: "1d7e2f23-b4c3-4ad7-8478-e770b68e6f11"},
			},
			expected: expected{
				Err:   auth.ErrFindAuthUserApiContractNotFound,
				Value: auth.DeleteAuthUserApiContractOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM public.auth_user_api_contract`).
					WithArgs("1d7e2f23-b4c3-4ad7-8478-e770b68e6f11").
					WillReturnResult(sqlmock.NewResult(0, 0))
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

			got, err := repo.DeleteAuthUserApiContract(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}

func TestAuthRepository_InsertAuthUserApiContractsFromRole(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.InsertAuthUserApiContractsFromRoleInput
	}

	type expected = testutil.Result[auth.InsertAuthUserApiContractsFromRoleOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should insert user api contracts from role contracts",
			args: args{
				ctx: context.Background(),
				input: auth.InsertAuthUserApiContractsFromRoleInput{
					UserId: "358cbaad-316e-4539-9949-2636cdbd7e89",
					RoleId: "0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
				},
			},
			expected: expected{
				Value: auth.InsertAuthUserApiContractsFromRoleOutput{InsertedCount: 3},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO public.auth_user_api_contract`).
					WithArgs(
						"358cbaad-316e-4539-9949-2636cdbd7e89",
						"0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
					).
					WillReturnResult(sqlmock.NewResult(0, 3))
			},
		},
		{
			name: "should return create grant error when insert fails",
			args: args{
				ctx: context.Background(),
				input: auth.InsertAuthUserApiContractsFromRoleInput{
					UserId: "358cbaad-316e-4539-9949-2636cdbd7e89",
					RoleId: "0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
				},
			},
			expected: expected{
				Err:   auth.ErrCreateAuthUserApiContract,
				Value: auth.InsertAuthUserApiContractsFromRoleOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO public.auth_user_api_contract`).
					WithArgs(
						"358cbaad-316e-4539-9949-2636cdbd7e89",
						"0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
					).
					WillReturnError(auth.ErrCreateAuthUserApiContract)
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

			got, err := repo.InsertAuthUserApiContractsFromRole(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
