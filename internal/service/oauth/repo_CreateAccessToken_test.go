package oauth_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/leporo/sqlf"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/dbtx"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
)

func TestOauthService_CreateAccessToken(t *testing.T) {
	var (
		expectedId           string = faker.UUIDHyphenated()
		expectedAccessToken  string = faker.Word()
		expectedRefreshToken string = faker.Word()
	)

	type args struct {
		ctx   context.Context
		input oauth.CreateAccessTokenInput
	}

	type expected = testutil.Result[oauth.CreateAccessTokenOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(sqlMock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "success create new token",
			args: args{
				ctx: context.Background(),
				input: oauth.CreateAccessTokenInput{
					Id:           expectedId,
					AccessToken:  expectedAccessToken,
					RefreshToken: expectedRefreshToken,
				},
			},
			expected: expected{
				Err: nil,
				Value: oauth.CreateAccessTokenOutput{
					Id: expectedId,
				},
			},
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(``).WithArgs(
					expectedId,
					expectedAccessToken,
					expectedRefreshToken,
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dbMock, sqlMock, _ := sqlmock.New()

			if tt.mock != nil {
				tt.mock(sqlMock)
			}

			dbTx := dbtx.DbTx{Db: dbMock}
			oauthWriteRepo := oauth.NewOAuthRepository(oauth.OAuthRepositoryOpts{
				DbTxInterface: &dbTx,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				Db:            dbMock,
			})
			got, err := oauthWriteRepo.CreateAccessToken(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
