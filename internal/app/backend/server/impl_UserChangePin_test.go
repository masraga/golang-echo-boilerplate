package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/app/backend/server"
	"github.com/masraga/golang-echo-boilerplate/internal/crypto"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestServer_UserChangePin(t *testing.T) {
	e := echo.New()

	var (
		userPhoneNo      = "encrypted-user-phone-no"
		oldPinHash       = "encrypted-old-pin"
		newPinHash       = "encrypted-new-pin"
		retypeNewPinHash = "encrypted-retype-new-pin"
		oldPin           = "1111"
		newPin           = "2222"
		retypeNewPin     = "2222"
	)

	type args struct {
		userId api.UserIdPathParameter
		input  api.UserChangePinRequest
	}

	type fields struct {
		AuthService   *auth.MockAuthServiceInterface
		CryptoService *crypto.MockCryptoServiceInterface
	}

	type expected = testutil.HttpResult

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "success with response OK",
			args: args{
				input: api.UserChangePinRequest{
					UserPhoneNo:  userPhoneNo,
					OldPin:       oldPinHash,
					NewPin:       newPinHash,
					RetypeNewPin: retypeNewPinHash,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				cryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				cryptoService.EXPECT().
					Decrypt(gomock.Any(), gomock.Any()).
					Return(crypto.DecryptOutput{Result: userPhoneNo}, nil)
				cryptoService.EXPECT().
					Decrypt(gomock.Any(), gomock.Any()).
					Return(crypto.DecryptOutput{Result: oldPin}, nil)
				cryptoService.EXPECT().
					Decrypt(gomock.Any(), gomock.Any()).
					Return(crypto.DecryptOutput{Result: newPin}, nil)
				cryptoService.EXPECT().
					Decrypt(gomock.Any(), gomock.Any()).
					Return(crypto.DecryptOutput{Result: retypeNewPin}, nil)

				authService := auth.NewMockAuthServiceInterface(ctrl)
				authService.EXPECT().
					UserChangePin(gomock.Any(), gomock.Any()).
					Return(auth.UserChangePinOutput{IsUpdate: true}, nil)

				tt.fields.AuthService = authService
				tt.fields.CryptoService = cryptoService

				tt.expected = testutil.HttpResult{
					Code: http.StatusOK,
					Body: `{"isUpdate": true}`,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			body, _ := json.Marshal(tt.args.input)
			req := httptest.NewRequest(
				"PUT",
				"/api/v1/auth/users/change-pin",
				strings.NewReader(string(body)),
			)

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			context := e.NewContext(req, rec)

			if tt.mock != nil {
				tt.mock(&tt, ctrl)
			}

			server := server.NewServer(server.ServerOpts{
				CryptoService: tt.fields.CryptoService,
				AuthService:   tt.fields.AuthService,
			})

			server.UserChangePin(context)

			testutil.RequireHttpResultJson(t, tt.expected, rec)
		})
	}
}
