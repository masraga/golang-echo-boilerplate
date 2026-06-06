package auth_test

import (
	"context"
	"testing"

	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestAuthService_AssignAuthUserRole(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.AssignAuthUserRoleInput
	}

	type expected = testutil.Result[auth.AssignAuthUserRoleOutput]

	type fields struct {
		AuthRepositoryReader *auth.MockAuthRepositoryReaderInterface
		AuthRepositoryWriter *auth.MockAuthRepositoryWriterInterface
	}

	type test struct {
		name     string
		args     args
		expected expected
		fields   fields
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should return role not found when role lookup fails",
			args: args{
				ctx: context.Background(),
				input: auth.AssignAuthUserRoleInput{
					UserId:    "358cbaad-316e-4539-9949-2636cdbd7e89",
					RoleId:    "0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
					CreatedBy: "admin",
				},
			},
			expected: expected{
				Err:   auth.ErrFindAuthRoleNotFound,
				Value: auth.AssignAuthUserRoleOutput{},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					GetAuthRole(tt.args.ctx, auth.GetAuthRoleInput{Id: tt.args.input.RoleId}).
					Return(auth.GetAuthRoleOutput{}, auth.ErrFindAuthRoleNotFound)

				tt.fields.AuthRepositoryReader = authRepoReader
			},
		},
		{
			name: "should replace user grants with role grants",
			args: args{
				ctx: context.Background(),
				input: auth.AssignAuthUserRoleInput{
					UserId:    "358cbaad-316e-4539-9949-2636cdbd7e89",
					RoleId:    "0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
					CreatedBy: "admin",
				},
			},
			expected: expected{
				Value: auth.AssignAuthUserRoleOutput{
					UserId:        "358cbaad-316e-4539-9949-2636cdbd7e89",
					RoleId:        "0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
					RoleName:      "finance-admin",
					GrantedCount:  3,
					UpdatedAtUtc0: 1798790400000,
					CreatedBy:     "admin",
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoReader := auth.NewMockAuthRepositoryReaderInterface(ctrl)
				authRepoReader.EXPECT().
					GetAuthRole(tt.args.ctx, auth.GetAuthRoleInput{Id: tt.args.input.RoleId}).
					Return(auth.GetAuthRoleOutput{
						Id:       tt.args.input.RoleId,
						RoleName: "finance-admin",
					}, nil)

				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(tt.args.ctx, nil).
					Return(tt.args.ctx, nil)
				authRepoWriter.EXPECT().
					AssignAuthUserRole(tt.args.ctx, auth.AssignAuthUserRoleInput{
						UserId:    tt.args.input.UserId,
						RoleId:    tt.args.input.RoleId,
						RoleName:  "finance-admin",
						CreatedBy: tt.args.input.CreatedBy,
					}).
					Return(auth.AssignAuthUserRoleOutput{
						UserId:        tt.args.input.UserId,
						RoleId:        tt.args.input.RoleId,
						RoleName:      "finance-admin",
						UpdatedAtUtc0: 1798790400000,
						CreatedBy:     tt.args.input.CreatedBy,
					}, nil)
				authRepoWriter.EXPECT().
					DeleteAuthUserApiContractsByUserId(tt.args.ctx, auth.DeleteAuthUserApiContractsByUserIdInput{
						UserId: tt.args.input.UserId,
					}).
					Return(auth.DeleteAuthUserApiContractsByUserIdOutput{DeletedCount: 1}, nil)
				authRepoWriter.EXPECT().
					InsertAuthUserApiContractsFromRole(tt.args.ctx, auth.InsertAuthUserApiContractsFromRoleInput{
						UserId:    tt.args.input.UserId,
						RoleId:    tt.args.input.RoleId,
						CreatedBy: tt.args.input.CreatedBy,
					}).
					Return(auth.InsertAuthUserApiContractsFromRoleOutput{InsertedCount: 3}, nil)
				authRepoWriter.EXPECT().
					CommitOrRollback(tt.args.ctx, nil).
					Return(nil)

				tt.fields.AuthRepositoryReader = authRepoReader
				tt.fields.AuthRepositoryWriter = authRepoWriter
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

			authService := auth.NewAuthService(auth.AuthServiceOpts{
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryReader: tt.fields.AuthRepositoryReader,
				AuthRepositoryWriter: tt.fields.AuthRepositoryWriter,
			})

			got, err := authService.AssignAuthUserRole(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}

func TestAuthService_DeleteAuthUserRole(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.DeleteAuthUserRoleInput
	}

	type expected = testutil.Result[auth.DeleteAuthUserRoleOutput]

	type fields struct {
		AuthRepositoryWriter *auth.MockAuthRepositoryWriterInterface
	}

	type test struct {
		name     string
		args     args
		expected expected
		fields   fields
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should clear role and hard delete user grants",
			args: args{
				ctx: context.Background(),
				input: auth.DeleteAuthUserRoleInput{
					UserId:    "358cbaad-316e-4539-9949-2636cdbd7e89",
					CreatedBy: "admin",
				},
			},
			expected: expected{
				Value: auth.DeleteAuthUserRoleOutput{IsSuccess: true},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				authRepoWriter := auth.NewMockAuthRepositoryWriterInterface(ctrl)
				authRepoWriter.EXPECT().
					Begin(tt.args.ctx, nil).
					Return(tt.args.ctx, nil)
				authRepoWriter.EXPECT().
					DeleteAuthUserRole(tt.args.ctx, tt.args.input).
					Return(auth.DeleteAuthUserRoleOutput{IsSuccess: true}, nil)
				authRepoWriter.EXPECT().
					DeleteAuthUserApiContractsByUserId(tt.args.ctx, auth.DeleteAuthUserApiContractsByUserIdInput{
						UserId: tt.args.input.UserId,
					}).
					Return(auth.DeleteAuthUserApiContractsByUserIdOutput{DeletedCount: 2}, nil)
				authRepoWriter.EXPECT().
					CommitOrRollback(tt.args.ctx, nil).
					Return(nil)

				tt.fields.AuthRepositoryWriter = authRepoWriter
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

			authService := auth.NewAuthService(auth.AuthServiceOpts{
				Err:                  ctxerr.NewCtxErr(ctxerr.CtxErrOpts{}),
				AuthRepositoryWriter: tt.fields.AuthRepositoryWriter,
			})

			got, err := authService.DeleteAuthUserRole(tt.args.ctx, tt.args.input)
			testutil.RequireResult(t, err, tt.expected, got)
		})
	}
}
