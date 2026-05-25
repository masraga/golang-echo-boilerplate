package auth_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/dbtx"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
)

func TestAuthRepository_FindOTP(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.FindOTPInput
	}

	type expected = testutil.Result[auth.FindOTPOutput]

	type test struct {
		name     string
		expected expected
		args     args
		mock     func(mock sqlmock.Sqlmock)
	}

	var (
		expectedId      string = faker.UUIDHyphenated()
		expectedOtpCode string = "123456"
	)

	tests := []test{
		{
			name: "Should not found OTP",
			args: args{
				ctx: context.Background(),
				input: auth.FindOTPInput{
					UserId: "user-id",
				},
			},
			expected: expected{
				Err:   auth.ErrFindOTPNotFound,
				Value: auth.FindOTPOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(``).
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(auth.ErrFindOTPNotFound)
			},
		},
		{
			name: "Should success find OTP",
			args: args{
				ctx: context.Background(),
				input: auth.FindOTPInput{
					UserId: "user-id",
				},
			},
			expected: expected{
				Err:   nil,
				Value: auth.FindOTPOutput{Id: expectedId, OtpCode: expectedOtpCode},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "otp_code"}).
							AddRow(expectedId, expectedOtpCode),
					)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer dbMock.Close()

			if tt.mock != nil {
				tt.mock(sqlMock)
			}

			dbTx := &dbtx.DbTx{Db: dbMock}
			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: dbTx,
				Sql:           sqlf.PostgreSQL,
				Db:            dbMock,
			})

			res, err := repo.FindOTP(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
