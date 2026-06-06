package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/leporo/sqlf"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/dbtx"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
)

func TestAuthService_IntegrationValidateJwtToken(t *testing.T) {
	var (
		jwtSecret             = auth.JwtSecretType(faker.WORD)
		expectedUserId        = faker.UUIDHyphenated()
		expectedExpiredAtUtc0 = time.Now().Add(time.Hour).UnixMilli()
		expectedIssuerAtUtc0  = time.Now().UnixMilli()
	)

	type args struct {
		ctx   context.Context
		input auth.ValidateJwtTokenInput
	}

	type expected = testutil.Result[auth.ValidateJwtTokenOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(tt *test, authService *auth.AuthService, sqlMock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "Should valid token",
			args: args{
				ctx: context.Background(),
			},
			expected: expected{
				Err: nil,
				Value: auth.ValidateJwtTokenOutput{
					UserId:        expectedUserId,
					ExpiredAtUtc0: expectedExpiredAtUtc0,
					IssuerAtUtc0:  expectedIssuerAtUtc0,
				},
			},
			mock: func(tt *test, authService *auth.AuthService, sqlMock sqlmock.Sqlmock) {
				tokenOutput, _ := authService.CreateJWTToken(tt.args.ctx, auth.CreateJWTTokenInput{
					ExpiredAtUtc0: expectedExpiredAtUtc0,
					IssuerAtUtc0:  expectedIssuerAtUtc0,
					UserId:        expectedUserId,
				})
				tt.args.input.Token = tokenOutput.Token

				sqlMock.ExpectQuery(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "expired_at_utc0", "is_active"}).
						AddRow(tokenOutput.Token, expectedUserId, expectedExpiredAtUtc0, true))
			},
		},
		{
			name: "Should invalid token",
			args: args{
				ctx: context.Background(),
			},
			expected: expected{
				Err: auth.ErrAuthTokenExpired,
				Value: auth.ValidateJwtTokenOutput{
					UserId:        expectedUserId,
					ExpiredAtUtc0: time.Now().Add(-time.Hour).UnixMilli(),
					IssuerAtUtc0:  time.Now().Add(-time.Hour).UnixMilli(),
				},
			},
			mock: func(tt *test, authService *auth.AuthService, sqlMock sqlmock.Sqlmock) {
				tokenOutput, _ := authService.CreateJWTToken(tt.args.ctx, auth.CreateJWTTokenInput{
					ExpiredAtUtc0: time.Now().Add(-time.Hour).UnixMilli(),
					IssuerAtUtc0:  time.Now().Add(-time.Hour).UnixMilli(),
					UserId:        expectedUserId,
				})
				tt.args.input.Token = tokenOutput.Token
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

			dbTx := &dbtx.DbTx{Db: dbMock}
			authRepository := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: dbTx,
				Db:            dbMock,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})
			authService := auth.NewAuthService(auth.AuthServiceOpts{
				JwtSecret:            jwtSecret,
				JwtExpiration:        auth.JwtExpirationType(60),
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader: authRepository,
			})

			if tt.mock != nil {
				tt.mock(&tt, authService, sqlMock)
			}

			got, err := authService.ValidateJwtToken(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
