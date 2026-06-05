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

func TestAuthRepository_UpdateFirebaseId(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.UpdateFirebaseIdInput
	}

	type expected = testutil.Result[auth.UpdateFirebaseIdOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should fail to update firebase id",
			args: args{
				ctx: context.Background(),
				input: auth.UpdateFirebaseIdInput{
					UserId:     "358cbaad-316e-4539-9949-2636cdbd7e89",
					FirebaseId: "fcm-registration-token",
				},
			},
			expected: expected{
				Err: auth.ErrUpdateFirebaseId,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrUpdateFirebaseId)
			},
		},
		{
			name: "should update firebase id",
			args: args{
				ctx: context.Background(),
				input: auth.UpdateFirebaseIdInput{
					UserId:     "358cbaad-316e-4539-9949-2636cdbd7e89",
					FirebaseId: "fcm-registration-token",
				},
			},
			expected: expected{
				Value: auth.UpdateFirebaseIdOutput{
					UserId:     "358cbaad-316e-4539-9949-2636cdbd7e89",
					FirebaseId: "fcm-registration-token",
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

			if tt.mock != nil {
				tt.mock(sqlMock)
			}

			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: &dbtx.DbTx{Db: db},
				Db:            db,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})
			res, err := repo.UpdateFirebaseId(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
