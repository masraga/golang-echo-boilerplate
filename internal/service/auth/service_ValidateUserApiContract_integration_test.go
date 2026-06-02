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

func TestAuthService_IntegrationValidateUserApiContract(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.ValidateUserApiContractInput
	}

	type fields struct {
		AuthAccessBootstrapUserId auth.AuthAccessBootstrapUserIdType
	}

	type expected = testutil.Result[auth.ValidateUserApiContractOutput]

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should bypass database grant lookup for bootstrap user",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateUserApiContractInput{
					UserId:         "bootstrap-user-id",
					EndpointPath:   "/api/v1/crypto/encrypt",
					EndpointMethod: "post",
				},
			},
			fields: fields{
				AuthAccessBootstrapUserId: "bootstrap-user-id",
			},
			expected: expected{
				Value: auth.ValidateUserApiContractOutput{IsAllowed: true},
			},
		},
		{
			name: "should validate database grant for non bootstrap user",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateUserApiContractInput{
					UserId:         "user-id",
					EndpointPath:   "api/v1/crypto/encrypt",
					EndpointMethod: "POST",
				},
			},
			fields: fields{
				AuthAccessBootstrapUserId: "bootstrap-user-id",
			},
			expected: expected{
				Value: auth.ValidateUserApiContractOutput{IsAllowed: true},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(``).
					WithArgs("user-id", "/api/v1/crypto/encrypt", "post", true, true).
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

			authRepository := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: &dbtx.DbTx{Db: dbMock},
				Db:            dbMock,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})
			authService := auth.NewAuthService(auth.AuthServiceOpts{
				AuthAccessBootstrapUserId: tt.fields.AuthAccessBootstrapUserId,
				Err:                       ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader:      authRepository,
			})

			got, err := authService.ValidateUserApiContract(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
