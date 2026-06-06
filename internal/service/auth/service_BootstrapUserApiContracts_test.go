package auth_test

import (
	"context"
	"testing"

	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestAuthService_BootstrapUserApiContracts(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.BootstrapUserApiContractsInput
	}

	type fields struct {
		AuthRepositoryWriter *auth.MockAuthRepositoryWriterInterface
	}

	type expected = testutil.Result[auth.BootstrapUserApiContractsOutput]

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should bootstrap user api contracts",
			args: args{
				ctx:   context.Background(),
				input: auth.BootstrapUserApiContractsInput{UserId: "user-id"},
			},
			expected: expected{
				Value: auth.BootstrapUserApiContractsOutput{InsertedCount: 3},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				writer := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				writer.EXPECT().
					BootstrapUserApiContracts(gomock.Any(), tt.args.input).
					Return(auth.BootstrapUserApiContractsOutput{InsertedCount: 3}, nil)
				tt.fields.AuthRepositoryWriter = writer
			},
		},
		{
			name: "should return bootstrap error",
			args: args{
				ctx:   context.Background(),
				input: auth.BootstrapUserApiContractsInput{UserId: "user-id"},
			},
			expected: expected{
				Err:   auth.ErrBootstrapUserApiContract,
				Value: auth.BootstrapUserApiContractsOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				writer := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				writer.EXPECT().
					BootstrapUserApiContracts(gomock.Any(), tt.args.input).
					Return(auth.BootstrapUserApiContractsOutput{}, auth.ErrBootstrapUserApiContract)
				tt.fields.AuthRepositoryWriter = writer
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.mock != nil {
				tt.mock(&tt, ctrl)
			}

			svc := auth.NewAuthService(auth.AuthServiceOpts{
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryWriter: tt.fields.AuthRepositoryWriter,
			})

			got, err := svc.BootstrapUserApiContracts(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
