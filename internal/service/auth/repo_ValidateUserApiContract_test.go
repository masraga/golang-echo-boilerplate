package auth_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/dbtx"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
)

func TestAuthRepository_ValidateUserApiContract(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.ValidateUserApiContractInput
	}

	type expected = testutil.Result[auth.ValidateUserApiContractOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should return forbidden when grant not found",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateUserApiContractInput{
					UserId:         "358cbaad-316e-4539-9949-2636cdbd7e89",
					EndpointPath:   "/api/v1/crypto/encrypt",
					EndpointMethod: "post",
				},
			},
			expected: expected{
				Err:   auth.ErrUserApiContractForbidden,
				Value: auth.ValidateUserApiContractOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "should validate active grant",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateUserApiContractInput{
					UserId:         "358cbaad-316e-4539-9949-2636cdbd7e89",
					EndpointPath:   "/api/v1/crypto/encrypt",
					EndpointMethod: "post",
				},
			},
			expected: expected{
				Value: auth.ValidateUserApiContractOutput{IsAllowed: true},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("grant-id"))
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

			got, err := repo.ValidateUserApiContract(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
