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

func TestAuthRepository_DeleteAllUerOTP(t *testing.T) {
	var (
		expectedUserId string = faker.UUIDHyphenated()
	)

	type args struct {
		ctx   context.Context
		input auth.DeleteAllUserOTPInput
	}

	type expected = testutil.Result[auth.DeleteAllUserOTPOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should fail delete all user otp",
			args: args{
				ctx: context.Background(),
				input: auth.DeleteAllUserOTPInput{
					UserId: expectedUserId,
				},
			},
			expected: expected{
				Err:   auth.ErrFailedDeleteUserOTP,
				Value: auth.DeleteAllUserOTPOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrFailedDeleteUserOTP)
			},
		},
		{
			name: "should success delete all user otp",
			args: args{
				ctx: context.Background(),
				input: auth.DeleteAllUserOTPInput{
					UserId: expectedUserId,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.DeleteAllUserOTPOutput{
					IsSuccess: true,
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mock(mock)

			dbtx := &dbtx.DbTx{
				Db: db,
			}
			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: dbtx,
				Sql:           sqlf.PostgreSQL,
				Db:            db,
			})
			res, err := repo.DeleteAllUserOTP(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
