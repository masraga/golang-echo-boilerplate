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
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/app/backend/server"
	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	"github.com/masraga/kerp-api/internal/util/pointer"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestServer_AuthValidatePin(t *testing.T) {
	e := echo.New()

	var (
		encryptedPhoneNo = "encrypted-phone-no"
		phoneNo          = "081234567890"
		pinCode          = "123456"
		token            = "auth-token"
		userId           = uuid.MustParse("f7fa6a88-3af1-4d5f-b3b2-8f4cb0fd5784")
	)

	type args struct {
		input api.AuthValidatePinRequest
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
		mock     func(ctx echo.Context, tt *test, ctrl *gomock.Controller)
	}

	tests := []test{
		{
			name: "test success endpoint",
			args: args{
				input: api.AuthValidatePinRequest{
					PhoneNo:   encryptedPhoneNo,
					Pin:       pinCode,
					RetypePin: pointer.String(pinCode),
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(ctx.Request().Context(), crypto.DecryptInput{
						HashCode: tt.args.input.PhoneNo,
					}).
					Return(crypto.DecryptOutput{Result: phoneNo}, nil)

				mockAuthService := auth.NewMockAuthServiceInterface(ctrl)
				mockAuthService.EXPECT().
					AuthValidatePin(ctx.Request().Context(), auth.AuthValidatePinInput{
						PhoneNo:       phoneNo,
						PinCode:       tt.args.input.Pin,
						RetypePinCode: tt.args.input.RetypePin,
					}).
					Return(auth.AuthValidatePinOutput{
						IsValid: true,
						Token:   token,
						UserId:  userId.String(),
					}, nil)

				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService

				result, _ := json.Marshal(api.AuthValidatePinResponse{
					AuthToken: pointer.String(token),
					UserId:    &userId,
				})
				tt.expected.httpResult = &testutil.HttpResult{
					Code: http.StatusOK,
					Body: string(result),
				}
			},
		},
		{
			name: "test failed when decrypt phone number failed",
			args: args{
				input: api.AuthValidatePinRequest{
					PhoneNo:   encryptedPhoneNo,
					Pin:       pinCode,
					RetypePin: pointer.String(pinCode),
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				errDecrypt := errors.New("failed decrypt phone number")
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(ctx.Request().Context(), crypto.DecryptInput{
						HashCode: tt.args.input.PhoneNo,
					}).
					Return(crypto.DecryptOutput{}, errDecrypt)

				tt.fields.CryptoService = mockCryptoService
				tt.expected.httpResult = &testutil.HttpResult{
					Code: http.StatusInternalServerError,
					Body: `{"error":"unknown error"}`,
				}
			},
		},
		{
			name: "test failed when validate pin failed",
			args: args{
				input: api.AuthValidatePinRequest{
					PhoneNo:   encryptedPhoneNo,
					Pin:       pinCode,
					RetypePin: pointer.String(pinCode),
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(ctx.Request().Context(), crypto.DecryptInput{
						HashCode: tt.args.input.PhoneNo,
					}).
					Return(crypto.DecryptOutput{Result: phoneNo}, nil)

				mockAuthService := auth.NewMockAuthServiceInterface(ctrl)
				mockAuthService.EXPECT().
					AuthValidatePin(ctx.Request().Context(), auth.AuthValidatePinInput{
						PhoneNo:       phoneNo,
						PinCode:       tt.args.input.Pin,
						RetypePinCode: tt.args.input.RetypePin,
					}).
					Return(auth.AuthValidatePinOutput{}, auth.ErrPinCodeNotMatch)

				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService
				tt.expected.httpResult = &testutil.HttpResult{
					Code: http.StatusInternalServerError,
					Body: `{"error":"unknown error"}`,
				}
			},
		},
		{
			name: "test failed when otp is not validated",
			args: args{
				input: api.AuthValidatePinRequest{
					PhoneNo: encryptedPhoneNo,
					Pin:     pinCode,
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(ctx.Request().Context(), crypto.DecryptInput{
						HashCode: tt.args.input.PhoneNo,
					}).
					Return(crypto.DecryptOutput{Result: phoneNo}, nil)

				mockAuthService := auth.NewMockAuthServiceInterface(ctrl)
				mockAuthService.EXPECT().
					AuthValidatePin(ctx.Request().Context(), auth.AuthValidatePinInput{
						PhoneNo: phoneNo,
						PinCode: tt.args.input.Pin,
					}).
					Return(auth.AuthValidatePinOutput{}, auth.ErrOtpValidationRequired)

				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService
				tt.expected.httpResult = &testutil.HttpResult{
					Code: http.StatusBadRequest,
					Body: `{"error":"must validate otp before validating pin"}`,
				}
			},
		},
		{
			name: "test failed when service returns invalid user id",
			args: args{
				input: api.AuthValidatePinRequest{
					PhoneNo:   encryptedPhoneNo,
					Pin:       pinCode,
					RetypePin: pointer.String(pinCode),
				},
			},
			mock: func(ctx echo.Context, tt *test, ctrl *gomock.Controller) {
				mockCryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
				mockCryptoService.EXPECT().
					Decrypt(ctx.Request().Context(), crypto.DecryptInput{
						HashCode: tt.args.input.PhoneNo,
					}).
					Return(crypto.DecryptOutput{Result: phoneNo}, nil)

				mockAuthService := auth.NewMockAuthServiceInterface(ctrl)
				mockAuthService.EXPECT().
					AuthValidatePin(ctx.Request().Context(), auth.AuthValidatePinInput{
						PhoneNo:       phoneNo,
						PinCode:       tt.args.input.Pin,
						RetypePinCode: tt.args.input.RetypePin,
					}).
					Return(auth.AuthValidatePinOutput{
						IsValid: true,
						Token:   token,
						UserId:  "invalid-user-id",
					}, nil)

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

			body, _ := json.Marshal(tt.args.input)
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/auth/pin/validate",
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

			err := svc.AuthValidatePin(context)
			require.NoError(t, err)
			testutil.RequireHttpResultJson(t, *tt.expected.httpResult, rec)
		})
	}
}
