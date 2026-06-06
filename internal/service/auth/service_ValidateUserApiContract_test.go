package auth_test

import (
	"context"
	"testing"

	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestAuthService_ValidateUserApiContract(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.ValidateUserApiContractInput
	}

	type fields struct {
		AuthAccessBootstrapUserId auth.AuthAccessBootstrapUserIdType
		AuthRepositoryReader      *auth.MockAuthRepositoryReaderInterface
	}

	type expected = testutil.Result[auth.ValidateUserApiContractOutput]

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should bypass repository validation for bootstrap user",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateUserApiContractInput{
					UserId:         "bootstrap-user-id",
					EndpointPath:   "/api/v1/crypto/encrypt",
					EndpointMethod: "post",
				},
			},
			fields: fields{
				AuthAccessBootstrapUserId: "bootstrap-user-id",
			},
			expected: expected{
				Value: auth.ValidateUserApiContractOutput{IsAllowed: true},
			},
		},
		{
			name: "should normalize path and method then validate access",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateUserApiContractInput{
					UserId:         "user-id",
					EndpointPath:   "api/v1/crypto/encrypt",
					EndpointMethod: "POST",
				},
			},
			expected: expected{
				Value: auth.ValidateUserApiContractOutput{IsAllowed: true},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				reader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				reader.EXPECT().
					ValidateUserApiContract(gomock.Any(), auth.ValidateUserApiContractInput{
						UserId:         "user-id",
						EndpointPath:   "/api/v1/crypto/encrypt",
						EndpointMethod: "post",
					}).
					Return(auth.ValidateUserApiContractOutput{IsAllowed: true}, nil)
				tt.fields.AuthRepositoryReader = reader
			},
		},
		{
			name: "should return repository forbidden error",
			args: args{
				ctx: context.Background(),
				input: auth.ValidateUserApiContractInput{
					UserId:         "user-id",
					EndpointPath:   "/api/v1/crypto/encrypt",
					EndpointMethod: "post",
				},
			},
			expected: expected{
				Err:   auth.ErrUserApiContractForbidden,
				Value: auth.ValidateUserApiContractOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				reader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				reader.EXPECT().
					ValidateUserApiContract(gomock.Any(), auth.ValidateUserApiContractInput{
						UserId:         "user-id",
						EndpointPath:   "/api/v1/crypto/encrypt",
						EndpointMethod: "post",
					}).
					Return(auth.ValidateUserApiContractOutput{}, auth.ErrUserApiContractForbidden)
				tt.fields.AuthRepositoryReader = reader
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
				AuthAccessBootstrapUserId: tt.fields.AuthAccessBootstrapUserId,
				Err:                       ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader:      tt.fields.AuthRepositoryReader,
			})

			got, err := svc.ValidateUserApiContract(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
