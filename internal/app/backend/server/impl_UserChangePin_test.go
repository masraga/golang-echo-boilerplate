package server_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/app/backend/server"
	"github.com/masraga/golang-echo-boilerplate/internal/crypto"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestServer_UserChangePin(t *testing.T) {
	e := echo.New()

	var (
		userId           = uuid.MustParse("d290f1ee-6c54-4b01-90e6-d701748f0851")
		oldPinHash       = "encrypted-old-pin"
		newPinHash       = "encrypted-new-pin"
		retypeNewPinHash = "encrypted-retype-new-pin"
		oldPin           = "1111"
		newPin           = "2222"
		retypeNewPin     = "2222"
	)

	type args struct {
		input api.UserChangePinRequest
	}

	type fields struct {
		AuthService   *auth.MockAuthServiceInterface
		CryptoService *crypto.MockCryptoServiceInterface
	}

	type expected struct {
		httpResult *testutil.HttpResult
	}

	type test struct {
		name     string
		args     args
		fields   fields
		expected expected
		mock     func(tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "success",
			args: args{
				input: api.UserChangePinRequest{
					OldPin:       oldPinHash,
					NewPin:       newPinHash,
					RetypeNewPin: retypeNewPinHash,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(gomock.Any(), crypto.DecryptInput{HashCode: tt.args.input.OldPin}).
					Return(crypto.DecryptOutput{Result: oldPin}, nil)
				mockCryptoService.EXPECT().
					Decrypt(gomock.Any(), crypto.DecryptInput{HashCode: tt.args.input.NewPin}).
					Return(crypto.DecryptOutput{Result: newPin}, nil)
				mockCryptoService.EXPECT().
					Decrypt(gomock.Any(), crypto.DecryptInput{HashCode: tt.args.input.RetypeNewPin}).
					Return(crypto.DecryptOutput{Result: retypeNewPin}, nil)

				mockAuthService := auth.NewMockAuthServiceInterface(ctrl)
				mockAuthService.EXPECT().
					UserChangePin(gomock.Any(), auth.UserChangePinInput{
						UserId:       userId.String(),
						OldPin:       oldPin,
						NewPin:       newPin,
						RetypeNewPin: retypeNewPin,
					}).
					Return(auth.UserChangePinOutput{IsUpdate: true}, nil)

				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService

				result, _ := json.Marshal(api.UserChangePinResponse{IsUpdate: true})
				tt.expected.httpResult = &testutil.HttpResult{
					Code: http.StatusOK,
					Body: string(result),
				}
			},
		},
		{
			name: "decrypt old pin failure",
			args: args{
				input: api.UserChangePinRequest{
					OldPin:       oldPinHash,
					NewPin:       newPinHash,
					RetypeNewPin: retypeNewPinHash,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(gomock.Any(), crypto.DecryptInput{HashCode: tt.args.input.OldPin}).
					Return(crypto.DecryptOutput{}, errors.New("decrypt failed"))

				tt.fields.CryptoService = mockCryptoService
				tt.expected.httpResult = &testutil.HttpResult{
					Code: http.StatusInternalServerError,
					Body: `{"error":"unknown error"}`,
				}
			},
		},
		{
			name: "service returns pin mismatch",
			args: args{
				input: api.UserChangePinRequest{
					OldPin:       oldPinHash,
					NewPin:       newPinHash,
					RetypeNewPin: retypeNewPinHash,
				},
			},
			mock: func(tt *test, ctrl *gomock.Controller) {
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(gomock.Any(), crypto.DecryptInput{HashCode: tt.args.input.OldPin}).
					Return(crypto.DecryptOutput{Result: oldPin}, nil)
				mockCryptoService.EXPECT().
					Decrypt(gomock.Any(), crypto.DecryptInput{HashCode: tt.args.input.NewPin}).
					Return(crypto.DecryptOutput{Result: newPin}, nil)
				mockCryptoService.EXPECT().
					Decrypt(gomock.Any(), crypto.DecryptInput{HashCode: tt.args.input.RetypeNewPin}).
					Return(crypto.DecryptOutput{Result: retypeNewPin}, nil)

				mockAuthService := auth.NewMockAuthServiceInterface(ctrl)
				mockAuthService.EXPECT().
					UserChangePin(gomock.Any(), auth.UserChangePinInput{
						UserId:       userId.String(),
						OldPin:       oldPin,
						NewPin:       newPin,
						RetypeNewPin: retypeNewPin,
					}).
					Return(auth.UserChangePinOutput{}, auth.ErrPinCodeNotMatch)

				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService
				tt.expected.httpResult = &testutil.HttpResult{
					Code: http.StatusInternalServerError,
					Body: `{"error":"unknown error"}`,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			body, err := json.Marshal(tt.args.input)
			require.NoError(t, err)

			req := httptest.NewRequest(
				http.MethodPut,
				"/api/v1/auth/users/"+userId.String()+"/change-pin",
				strings.NewReader(string(body)),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			context := e.NewContext(req, rec)
			if tt.mock != nil {
				tt.mock(&tt, ctrl)
			}

			svc := server.NewServer(server.ServerOpts{
				AuthService:   tt.fields.AuthService,
				CryptoService: tt.fields.CryptoService,
			})

			err = svc.UserChangePin(context, api.UserIdPathParameter(userId.String()))
			require.NoError(t, err)
			testutil.RequireHttpResultJson(t, *tt.expected.httpResult, rec)
		})
	}
}
