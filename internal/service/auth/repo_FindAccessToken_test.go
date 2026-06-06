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

func TestAuthRepository_FindAccessToken(t *testing.T) {
	var (
		expectedToken     = "access-token"
		expectedUserId    = faker.UUIDHyphenated()
		expectedExpiredAt = time.Now().Add(time.Duration(auth.JwtTokenExpiredDuration) * time.Minute).UnixMilli()
	)

	type args struct {
		ctx   context.Context
		input auth.FindAccessTokenInput
	}

	type expected = testutil.Result[auth.FindAccessTokenOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should failed when access token not found",
			args: args{
				ctx: context.Background(),
				input: auth.FindAccessTokenInput{
					Token:  expectedToken,
					UserId: expectedUserId,
				},
			},
			expected: expected{
				Err:   auth.ErrFindAccessTokenNotFound,
				Value: auth.FindAccessTokenOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrFindAccessTokenNotFound)
			},
		},
		{
			name: "should success find access token",
			args: args{
				ctx: context.Background(),
				input: auth.FindAccessTokenInput{
					Token:  expectedToken,
					UserId: expectedUserId,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.FindAccessTokenOutput{
					Token:         expectedToken,
					UserId:        expectedUserId,
					ExpiredAtUtc0: expectedExpiredAt,
					IsActive:      true,
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "expired_at_utc0", "is_active"}).
						AddRow(expectedToken, expectedUserId, expectedExpiredAt, true))
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

			res, err := repo.FindAccessToken(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
