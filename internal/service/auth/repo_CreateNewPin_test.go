package auth_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/leporo/sqlf"
	"github.com/masraga/kerp-api/internal/ctxerr"
	"github.com/masraga/kerp-api/internal/dbtx"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
)

func TestAuthRepository_CreateNewPin(t *testing.T) {
	var (
		expectedUserId  = faker.UUIDHyphenated()
		expectedPinCode = "123456"
	)

	type args struct {
		ctx   context.Context
		input auth.CreateNewPinInput
	}

	type expected = testutil.Result[auth.CreateNewPinOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "should failed to create new pin",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewPinInput{
					UserId:  expectedUserId,
					PinCode: expectedPinCode,
				},
			},
			expected: expected{
				Err:   auth.ErrCreateNewPin,
				Value: auth.CreateNewPinOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrCreateNewPin)
			},
		},
		{
			name: "should success to create new pin",
			args: args{
				ctx: context.Background(),
				input: auth.CreateNewPinInput{
					UserId:  expectedUserId,
					PinCode: expectedPinCode,
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.CreateNewPinOutput{
					UserId:  expectedUserId,
					PinCode: expectedPinCode,
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

			res, err := repo.CreateNewPin(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
