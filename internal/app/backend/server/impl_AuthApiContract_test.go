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

func TestServer_CreateAuthApiContract(t *testing.T) {
	type args struct {
		input api.CreateAuthApiContractRequest
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
			name: "should create auth api contract",
			args: args{
				input: api.CreateAuthApiContractRequest{
					Id:             "CryptoEncryptText",
					EndpointPath:   "/api/v1/crypto/encrypt",
					EndpointMethod: "post",
					Description:    "Encrypt plain text",
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				authAccessService := auth.NewMockAuthServiceInterface(ctrl)
				authAccessService.EXPECT().
					CreateAuthApiContract(ctx.Request().Context(), auth.CreateAuthApiContractInput{
						Id:             tt.args.input.Id,
						EndpointPath:   tt.args.input.EndpointPath,
						EndpointMethod: tt.args.input.EndpointMethod,
						Description:    tt.args.input.Description,
					}).
					Return(auth.CreateAuthApiContractOutput{
						Id:             tt.args.input.Id,
						EndpointPath:   tt.args.input.EndpointPath,
						EndpointMethod: tt.args.input.EndpointMethod,
						Description:    tt.args.input.Description,
						CreatedAtUtc0:  1798790400000,
						UpdatedAtUtc0:  1798790400000,
						IsActive:       true,
					}, nil)
				tt.fields.AuthService = authAccessService

				result, _ := json.Marshal(api.AuthApiContract{
					Id:             tt.args.input.Id,
					EndpointPath:   tt.args.input.EndpointPath,
					EndpointMethod: tt.args.input.EndpointMethod,
					Description:    tt.args.input.Description,
					CreatedAtUtc0:  1798790400000,
					UpdatedAtUtc0:  1798790400000,
					IsActive:       true,
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
				input: api.CreateAuthApiContractRequest{
					Id:             "CryptoEncryptText",
					EndpointPath:   "/api/v1/crypto/encrypt",
					EndpointMethod: "post",
					Description:    "Encrypt plain text",
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				authAccessService := auth.NewMockAuthServiceInterface(ctrl)
				authAccessService.EXPECT().
					CreateAuthApiContract(ctx.Request().Context(), gomock.Any()).
					Return(auth.CreateAuthApiContractOutput{}, auth.ErrCreateAuthApiContract)
				tt.fields.AuthService = authAccessService
				tt.expected.httpResult = testutil.HttpResult{
					Code: http.StatusBadRequest,
					Body: `{"error":"error to create auth api contract"}`,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			body, _ := json.Marshal(tt.args.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/api-contracts", strings.NewReader(string(body)))
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
			err := svc.CreateAuthApiContract(ctx)
			if err != nil {
				t.Fatal(err)
			}
			testutil.RequireHttpResultJson(t, tt.expected.httpResult, rec)
		})
	}
}
