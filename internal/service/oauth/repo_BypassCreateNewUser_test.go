package oauth_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/leporo/sqlf"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/dbtx"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/rs/zerolog"
)

func TestOAuthRepository_BypassCreateNewUser(t *testing.T) {
	var (
		expectedPhoneNo   string = "081234567890"
		expectedPin       string = "123456"
		expectedEmail     string = "admin@gmail.com"
		expectedCreatedBy string = oauth.DefaultCreatedBy
		expectedLoginMode string = "GOOGLE"
		expectedId        string = faker.UUIDHyphenated()
	)

	type args struct {
		ctx   context.Context
		input oauth.BypassCreateNewUserInput
	}

	type expected = testutil.Result[oauth.BypassCreateNewUserOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(sqlMock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "success bypass new user",
			args: args{
				ctx: context.Background(),
				input: oauth.BypassCreateNewUserInput{
					Id:        expectedId,
					PhoneNo:   expectedPhoneNo,
					Pin:       expectedPin,
					Email:     expectedEmail,
					CreatedBy: expectedCreatedBy,
					LoginMode: auth.LoginMode(expectedLoginMode),
				},
			},
			expected: expected{
				Err: nil,
				Value: oauth.BypassCreateNewUserOutput{
					Id: expectedId,
				},
			},
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectExec(``).
					WithArgs(
						expectedId,
						expectedPhoneNo,
						expectedPin,
						true,
						expectedEmail,
						expectedCreatedBy,
						expectedLoginMode,
					).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dbMock, sqlMock, _ := sqlmock.New()
			defer dbMock.Close()

			if tt.mock != nil {
				tt.mock(sqlMock)
			}

			dbtx := dbtx.DbTx{Db: dbMock}
			repo := oauth.NewOAuthRepository(oauth.OAuthRepositoryOpts{
				DbTxInterface: &dbtx,
				Sql:           sqlf.PostgreSQL,
				Db:            dbMock,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{Logger: zerolog.Nop()}),
			})
			got, err := repo.BypassCreateNewUser(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
