package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/app/backend/server"
	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestServer_RegisterPhoneNumber(t *testing.T) {
	e := echo.New()

	var (
		expectedId string = faker.Word()
	)

	type args struct {
		input api.CreateNewAccountRequest
	}

	type fields struct {
		AuthService   *auth.MockAuthServiceInterface
		CryptoService *crypto.MockCryptoServiceInterface
	}

	type test struct {
		name     string
		args     args
		fields   fields
		expected testutil.HttpResult
		mock     func(ctx echo.Context, tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "test success endpoint",
			args: args{
				input: api.CreateNewAccountRequest{
					PhoneNo: "081234567890",
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(ctx.Request().Context(), crypto.DecryptInput{
						HashCode: tt.args.input.PhoneNo,
					}).
					Return(crypto.DecryptOutput{Result: tt.args.input.PhoneNo}, nil)

				mockAuthService := auth.NewMockAuthServiceInterface(ctrl)
				mockAuthService.EXPECT().
					CreateNewAccount(ctx.Request().Context(), auth.CreateNewAccountInput{
						PhoneNo: tt.args.input.PhoneNo,
					}).
					Return(auth.CreateNewAccountOutput{Id: expectedId}, nil)
				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService

				result, _ := json.Marshal(api.CreateNewAccountResponse{
					Id: expectedId,
				})
				tt.expected = testutil.HttpResult{
					Code: http.StatusCreated,
					Body: string(result),
				}
			},
		},
		{
			name: "test failed when user already registered",
			args: args{
				input: api.CreateNewAccountRequest{
					PhoneNo: "081234567890",
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(ctx.Request().Context(), crypto.DecryptInput{
						HashCode: tt.args.input.PhoneNo,
					}).
					Return(crypto.DecryptOutput{Result: tt.args.input.PhoneNo}, nil)

				mockAuthService := auth.NewMockAuthServiceInterface(ctrl)
				mockAuthService.EXPECT().
					CreateNewAccount(ctx.Request().Context(), auth.CreateNewAccountInput{
						PhoneNo: tt.args.input.PhoneNo,
					}).
					Return(auth.CreateNewAccountOutput{}, auth.ErrDuplicateUser)
				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService

				result, _ := json.Marshal(api.ErrorBadRequest{Error: auth.ErrDuplicateUser.Error()})
				tt.expected = testutil.HttpResult{
					Code: http.StatusBadRequest,
					Body: string(result),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			body, _ := json.Marshal(tt.args.input)
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/auth/register/phone",
				strings.NewReader(string(body)),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			context := e.NewContext(req, rec)
			if tt.mock != nil {
				tt.mock(context, &tt, ctrl)
			}

			svc := server.NewServer(server.ServerOpts{
				AuthService:   tt.fields.AuthService,
				CryptoService: tt.fields.CryptoService,
			})
			svc.RegisterPhoneNumber(context)

			testutil.RequireHttpResultJson(t, tt.expected, rec)
		})
	}
}
