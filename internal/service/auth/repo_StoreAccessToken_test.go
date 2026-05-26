package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/dbtx"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
)

func TestAuthRepository_StoreAccessToken(t *testing.T) {
	var (
		expectedToken     = "access-token"
		expectedUserId    = faker.UUIDHyphenated()
		expectedExpiredAt = time.Now().Add(time.Duration(auth.JwtTokenExpiredDuration) * time.Minute).UnixMilli()
	)

	type args struct {
		ctx   context.Context
		input auth.StoreAccessTokenInput
	}

	type expected = testutil.Result[auth.StoreAccessTokenOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should failed to deactivate old access token",
			args: args{
				ctx: context.Background(),
				input: auth.StoreAccessTokenInput{
					Token:         expectedToken,
					UserId:        expectedUserId,
					ExpiredAtUtc0: expectedExpiredAt,
				},
			},
			expected: expected{
				Err:   auth.ErrStoreAccessToken,
				Value: auth.StoreAccessTokenOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrStoreAccessToken)
			},
		},
		{
			name: "should failed to insert access token",
			args: args{
				ctx: context.Background(),
				input: auth.StoreAccessTokenInput{
					Token:         expectedToken,
					UserId:        expectedUserId,
					ExpiredAtUtc0: expectedExpiredAt,
				},
			},
			expected: expected{
				Err:   auth.ErrStoreAccessToken,
				Value: auth.StoreAccessTokenOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrStoreAccessToken)
			},
		},
		{
			name: "should success to store access token",
			args: args{
				ctx: context.Background(),
				input: auth.StoreAccessTokenInput{
					Token:         expectedToken,
					UserId:        expectedUserId,
					ExpiredAtUtc0: expectedExpiredAt,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.StoreAccessTokenOutput{
					Token:         expectedToken,
					UserId:        expectedUserId,
					ExpiredAtUtc0: expectedExpiredAt,
					IsActive:      true,
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
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

			dbTx := &dbtx.DbTx{Db: dbMock}
			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: dbTx,
				Db:            dbMock,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})

			res, err := repo.StoreAccessToken(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
