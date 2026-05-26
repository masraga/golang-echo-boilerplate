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
	"github.com/masraga/kerp-api/internal/util/pointer"
)

func TestAuthRepository_CreateOTP(t *testing.T) {
	var (
		expectedOtpCode string = "123456"
	)
	type args struct {
		ctx   context.Context
		input auth.CreateOTPInput
	}

	type expected = testutil.Result[auth.CreateOTPOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should failed to add new otp",
			args: args{
				ctx: context.Background(),
				input: auth.CreateOTPInput{
					UserId:        faker.UUIDHyphenated(),
					Note:          pointer.String(faker.Word()),
					OtpCode:       expectedOtpCode,
					ExpiredAtUtc0: time.Now().Add(time.Duration(auth.OtpExpiredDuration) * time.Minute).UnixMilli(),
				},
			},
			expected: expected{
				Err:   auth.ErrCreateNewOTP,
				Value: auth.CreateOTPOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrCreateNewOTP)
			},
		},
		{
			name: "should success to add new otp",
			args: args{
				ctx: context.Background(),
				input: auth.CreateOTPInput{
					UserId:  faker.UUIDHyphenated(),
					Note:    pointer.String(faker.Word()),
					OtpCode: expectedOtpCode,
				},
			},
			expected: expected{
				Err:   nil,
				Value: auth.CreateOTPOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			dbTx := &dbtx.DbTx{Db: db}
			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: dbTx,
				Db:            db,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})

			res, err := repo.CreateOTP(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
