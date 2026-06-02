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

func TestServer_CreateAuthRoleContractApi(t *testing.T) {
	type args struct {
		roleId string
		input  api.CreateAuthRoleContractApiRequest
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
			name: "should create auth role contract api",
			args: args{
				roleId: "0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
				input: api.CreateAuthRoleContractApiRequest{
					AuthApiContractId: "CryptoEncryptText",
					CreatedBy:         "admin",
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				authService := auth.NewMockAuthServiceInterface(ctrl)
				authService.EXPECT().
					CreateAuthRoleContractApi(ctx.Request().Context(), auth.CreateAuthRoleContractApiInput{
						RoleId:            tt.args.roleId,
						AuthApiContractId: tt.args.input.AuthApiContractId,
						CreatedBy:         tt.args.input.CreatedBy,
					}).
					Return(auth.CreateAuthRoleContractApiOutput{
						Id:                "2c263a99-6c1e-45d1-9c0f-74088fcf50c7",
						RoleId:            tt.args.roleId,
						AuthApiContractId: tt.args.input.AuthApiContractId,
						CreatedAtUtc0:     1798790400000,
						UpdatedAtUtc0:     1798790400000,
						CreatedBy:         tt.args.input.CreatedBy,
						IsActive:          true,
					}, nil)
				tt.fields.AuthService = authService

				result, _ := json.Marshal(api.AuthRoleContractApi{
					Id:                "2c263a99-6c1e-45d1-9c0f-74088fcf50c7",
					RoleId:            tt.args.roleId,
					AuthApiContractId: tt.args.input.AuthApiContractId,
					CreatedAtUtc0:     1798790400000,
					UpdatedAtUtc0:     1798790400000,
					CreatedBy:         tt.args.input.CreatedBy,
					IsActive:          true,
				})
				tt.expected.httpResult = testutil.HttpResult{
					Code: http.StatusCreated,
					Body: string(result),
				}
			},
		},
		{
			name: "should return bad request when service fails",
			args: args{
				roleId: "0d8b2805-9b7f-47dd-8ec3-ec44105d37c7",
				input: api.CreateAuthRoleContractApiRequest{
					AuthApiContractId: "CryptoEncryptText",
					CreatedBy:         "admin",
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				authService := auth.NewMockAuthServiceInterface(ctrl)
				authService.EXPECT().
					CreateAuthRoleContractApi(ctx.Request().Context(), gomock.Any()).
					Return(auth.CreateAuthRoleContractApiOutput{}, auth.ErrCreateAuthRoleContractApi)
				tt.fields.AuthService = authService
				tt.expected.httpResult = testutil.HttpResult{
					Code: http.StatusBadRequest,
					Body: `{"error":"error to create auth role contract api"}`,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			body, _ := json.Marshal(tt.args.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/roles/"+tt.args.roleId+"/api-contracts", strings.NewReader(string(body)))
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
			err := svc.CreateAuthRoleContractApi(ctx, tt.args.roleId)
			if err != nil {
				t.Fatal(err)
			}
			testutil.RequireHttpResultJson(t, tt.expected.httpResult, rec)
		})
	}
}
