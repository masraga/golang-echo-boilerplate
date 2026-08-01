package auth_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/leporo/sqlf"
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/dbtx"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestAuthReposiory_FindAuthWithEmail(t *testing.T) {
	var (
		expectedId      string = faker.UUIDHyphenated()
		expectedEmail   string = faker.Word()
		expectedPhoneNo string = faker.Word()
	)

	type args struct {
		ctx   context.Context
		input auth.FindAuthWithEmailInput
	}

	type expected = testutil.Result[auth.FindAuthWithEmailOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(sqlMock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "success with return auth email",
			args: args{
				ctx: context.Background(),
				input: auth.FindAuthWithEmailInput{
					Email: expectedEmail,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.FindAuthWithEmailOutput{
					Id:      expectedId,
					Email:   expectedEmail,
					PhoneNo: expectedPhoneNo,
				},
			},
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(``).
					WithArgs(true, expectedEmail).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "phone_no", "email"}).
							AddRow(expectedId, expectedPhoneNo, expectedEmail),
					)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dbMock, sqlMock, err := sqlmock.New()
			require.NoError(t, err)

			if tt.mock != nil {
				tt.mock(sqlMock)
			}

			dbtx := dbtx.DbTx{Db: dbMock}
			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: &dbtx,
				Db:            dbMock,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{Logger: zerolog.Nop()}),
			})
			got, err := repo.FindAuthWithEmail(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
