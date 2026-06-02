package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/app/backend/server"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestServer_AssignAuthUserRole(t *testing.T) {
	type args struct {
		userId string
		input  api.AssignAuthUserRoleRequest
	}

	type fields struct {
		AuthService *auth.MockAuthServiceInterface
	}

	type expected struct {
		httpResult testutil.HttpResult
	}

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(ctx echo.Context, tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "should assign auth user role",
			args: args{
				userId: "358cbaad-316e-4539-9949-2636cdbd7e89",
				input: api.AssignAuthUserRoleRequest{
					RoleId:    "0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
					CreatedBy: "admin",
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				authService := auth.NewMockAuthServiceInterface(ctrl)
				authService.EXPECT().
					AssignAuthUserRole(ctx.Request().Context(), auth.AssignAuthUserRoleInput{
						UserId:    tt.args.userId,
						RoleId:    tt.args.input.RoleId,
						CreatedBy: tt.args.input.CreatedBy,
					}).
					Return(auth.AssignAuthUserRoleOutput{
						UserId:        tt.args.userId,
						RoleId:        tt.args.input.RoleId,
						RoleName:      "finance-admin",
						GrantedCount:  3,
						UpdatedAtUtc0: 1798790400000,
						CreatedBy:     tt.args.input.CreatedBy,
					}, nil)
				tt.fields.AuthService = authService

				result, _ := json.Marshal(api.AssignAuthUserRoleResponse{
					UserId:        tt.args.userId,
					RoleId:        tt.args.input.RoleId,
					RoleName:      "finance-admin",
					GrantedCount:  3,
					UpdatedAtUtc0: 1798790400000,
					CreatedBy:     tt.args.input.CreatedBy,
				})
				tt.expected.httpResult = testutil.HttpResult{
					Code: http.StatusOK,
					Body: string(result),
				}
			},
		},
		{
			name: "should return bad request when service fails",
			args: args{
				userId: "358cbaad-316e-4539-9949-2636cdbd7e89",
				input: api.AssignAuthUserRoleRequest{
					RoleId:    "0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
					CreatedBy: "admin",
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				authService := auth.NewMockAuthServiceInterface(ctrl)
				authService.EXPECT().
					AssignAuthUserRole(ctx.Request().Context(), gomock.Any()).
					Return(auth.AssignAuthUserRoleOutput{}, auth.ErrAssignAuthUserRole)
				tt.fields.AuthService = authService
				tt.expected.httpResult = testutil.HttpResult{
					Code: http.StatusBadRequest,
					Body: `{"error":"error to assign auth user role"}`,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			body, _ := json.Marshal(tt.args.input)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/auth/users/"+tt.args.userId+"/role", strings.NewReader(string(body)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			if tt.mock != nil {
				tt.mock(ctx, &tt, ctrl)
			}

			svc := server.NewServer(server.ServerOpts{
				AuthService: tt.fields.AuthService,
			})
			err := svc.AssignAuthUserRole(ctx, tt.args.userId)
			if err != nil {
				t.Fatal(err)
			}
			testutil.RequireHttpResultJson(t, tt.expected.httpResult, rec)
		})
	}
}
