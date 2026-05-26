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
	"github.com/masraga/kerp-api/internal/util/pointer"
)

func TestAuthRepository_FindAuth(t *testing.T) {
	var (
		expectedId      string  = faker.UUIDHyphenated()
		expectedPhoneNo string  = "081234567890"
		expectedPinCode *string = pointer.String("123456")
	)

	type args struct {
		ctx   context.Context
		input auth.FindAuthInput
	}

	type expected = testutil.Result[auth.FindAuthOutput]

	type test struct {
		name     string
		args     args
		expected expected
		mock     func(mock sqlmock.Sqlmock)
	}

	tests := []test{
		{
			name: "Test now result row",
			args: args{
				ctx: context.Background(),
				input: auth.FindAuthInput{
					PhoneNo: "0812343234",
				},
			},
			expected: expected{
				Err:   auth.ErrAuthNotFound,
				Value: auth.FindAuthOutput{},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(auth.ErrAuthNotFound)
			},
		},
		{
			name: "Test successfully with result row",
			args: args{
				ctx: context.Background(),
				input: auth.FindAuthInput{
					PhoneNo: "081234567890",
				},
			},
			expected: expected{
				Err: nil,
				Value: auth.FindAuthOutput{
					Id:      expectedId,
					PhoneNo: expectedPhoneNo,
					PinCode: expectedPinCode,
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(``).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(
						[]string{
							"id",
							"phone_no",
							"pin",
						},
					).AddRow(expectedId, expectedPhoneNo, expectedPinCode))
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

			dbtx := &dbtx.DbTx{Db: dbMock}
			repo := auth.NewAuthRepository(auth.AuthRepositoryOpts{
				DbTxInterface: dbtx,
				Db:            dbMock,
				Sql:           sqlf.PostgreSQL,
				Err:           ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
			})

			res, err := repo.FindAuth(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, res)
		})
	}
}
