package auth_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/dbtx"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
)

func TestAuthRepository_UpdateOtpValidity(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.UpdateOtpValidityInput
	}

	type expected = testutil.Result[auth.UpdateOtpValidityOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should fail to update otp validity",
			args: args{
				ctx: context.Background(),
				input: auth.UpdateOtpValidityInput{
					UserId:     "358cbaad-316e-4539-9949-2636cdbd7e89",
					IsOtpValid: true,
				},
			},
			expected: expected{Err: auth.ErrUpdateOtpValidity},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrUpdateOtpValidity)
			},
		},
		{
			name: "should update otp validity",
			args: args{
				ctx: context.Background(),
				input: auth.UpdateOtpValidityInput{
					UserId:     "358cbaad-316e-4539-9949-2636cdbd7e89",
					IsOtpValid: true,
				},
			},
			expected: expected{
				Value: auth.UpdateOtpValidityOutput{
					UserId:     "358cbaad-316e-4539-9949-2636cdbd7e89",
					IsOtpValid: true,
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error init mock: %v", err)
			}
			defer db.Close()

			tt.mock(sqlMock)
			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: &dbtx.DbTx{Db: db},
				Db:            db,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})

			result, err := repo.UpdateOtpValidity(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, result)
		})
	}
}
