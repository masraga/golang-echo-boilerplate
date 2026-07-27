package oauth_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/leporo/sqlf"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/dbtx"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
	"github.com/rs/zerolog"
)

func TestOauthService_GetAccessTokenById(t *testing.T) {
	var (
		expectedId                 string = faker.UUIDHyphenated()
		expectedUsername           string = faker.Email()
		expectedEmail              string = faker.Email()
		expectedAccessToken        string = faker.Word()
		expectedAccessRefreshToken string = faker.Word()
		expectedExpiryAtUtc0       int64  = time.Now().UnixMilli()
	)

	type args struct {
		ctx   context.Context
		input oauth.GetAccessTokenByIdInput
	}

	type expected = testutil.Result[oauth.GetAccessTokenByIdOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(sqlMock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "success with acces token",
			args: args{
				ctx: context.Background(),
				input: oauth.GetAccessTokenByIdInput{
					AccessToken: expectedAccessToken,
				},
			},
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(``).
					WithArgs(true, expectedAccessToken).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{
								"id",
								"username",
								"email",
								"access_token",
								"refresh_token",
								"expiry_at_utc0",
							},
						).AddRow(
							expectedId,
							expectedUsername,
							expectedEmail,
							expectedAccessToken,
							expectedAccessRefreshToken,
							expectedExpiryAtUtc0,
						),
					)
			},
			expected: expected{
				Err: nil,
				Value: oauth.GetAccessTokenByIdOutput{
					Id:           expectedId,
					Username:     pointer.String(expectedUsername),
					Email:        pointer.String(expectedEmail),
					AccessToken:  expectedAccessToken,
					RefreshToken: expectedAccessRefreshToken,
					ExpiryAtUtc0: pointer.Int64(expectedExpiryAtUtc0),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock, sqlMock, _ := sqlmock.New()
			defer dbMock.Close()

			if tt.mock != nil {
				tt.mock(sqlMock)
			}

			dbTx := dbtx.DbTx{Db: dbMock}
			repo := oauth.NewOAuthRepository(oauth.OAuthRepositoryOpts{
				DbTxInterface: &dbTx,
				Sql:           sqlf.PostgreSQL,
				Db:            dbMock,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{Logger: zerolog.Nop()}),
			})

			got, err := repo.GetAccessTokenById(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
