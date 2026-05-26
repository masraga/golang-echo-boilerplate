package server_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/app/backend/server"
	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/service/auth"
	"github.com/masraga/kerp-api/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestServer_VerifyNewAuthUserOTP(t *testing.T) {
	e := echo.New()

	var (
		encryptedPhoneNo = "encrypted-phone-no"
		phoneNo          = "081234567890"
		otpCode          = "123456"
		note             = "OTP verified"
	)

	type args struct {
		input api.VerifyOTPRequest
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
				input: api.VerifyOTPRequest{
					PhoneNo: encryptedPhoneNo,
					Otp:     otpCode,
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
					VerifyOtp(ctx.Request().Context(), auth.VerifyOtpInput{
						PhoneNo: phoneNo,
						OtpCode: tt.args.input.Otp,
					}).
					Return(auth.VerifyOtpOutput{
						IsValid: true,
						Note:    &note,
					}, nil)
				mockAuthService.EXPECT().
					VerifyUserAccount(ctx.Request().Context(), auth.VerifyUserAccountInput{
						PhoneNo: phoneNo,
					}).
					Return(auth.VerifyUserAccountOutput{
						PhoneNo:   phoneNo,
						IsNewUser: true,
					}, nil)
				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService

				result, _ := json.Marshal(api.VerifyOTPResponse200{
					IsValid:   true,
					PhoneNo:   phoneNo,
					Note:      &note,
					IsNewUser: true,
				})
				tt.expected = testutil.HttpResult{
					Code: http.StatusOK,
					Body: string(result),
				}
			},
		},
		{
			name: "test failed when decrypt phone number failed",
			args: args{
				input: api.VerifyOTPRequest{
					PhoneNo: encryptedPhoneNo,
					Otp:     otpCode,
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

				result, _ := json.Marshal(map[string]string{"error": "unknown error"})
				tt.expected = testutil.HttpResult{
					Code: http.StatusInternalServerError,
					Body: string(result),
				}
			},
		},
		{
			name: "test failed when verify otp failed",
			args: args{
				input: api.VerifyOTPRequest{
					PhoneNo: encryptedPhoneNo,
					Otp:     otpCode,
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
					VerifyOtp(ctx.Request().Context(), auth.VerifyOtpInput{
						PhoneNo: phoneNo,
						OtpCode: tt.args.input.Otp,
					}).
					Return(auth.VerifyOtpOutput{}, auth.ErrVerifyOtp)
				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService

				result, _ := json.Marshal(map[string]string{"error": "unknown error"})
				tt.expected = testutil.HttpResult{
					Code: http.StatusInternalServerError,
					Body: string(result),
				}
			},
		},
		{
			name: "test failed when verify user account failed",
			args: args{
				input: api.VerifyOTPRequest{
					PhoneNo: encryptedPhoneNo,
					Otp:     otpCode,
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
					VerifyOtp(ctx.Request().Context(), auth.VerifyOtpInput{
						PhoneNo: phoneNo,
						OtpCode: tt.args.input.Otp,
					}).
					Return(auth.VerifyOtpOutput{
						IsValid: true,
						Note:    &note,
					}, nil)
				mockAuthService.EXPECT().
					VerifyUserAccount(ctx.Request().Context(), auth.VerifyUserAccountInput{
						PhoneNo: phoneNo,
					}).
					Return(auth.VerifyUserAccountOutput{}, auth.ErrVerifyUserAccount)
				tt.fields.AuthService = mockAuthService
				tt.fields.CryptoService = mockCryptoService

				result, _ := json.Marshal(map[string]string{"error": "unknown error"})
				tt.expected = testutil.HttpResult{
					Code: http.StatusInternalServerError,
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
				"/api/v1/auth/otp/verify",
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
			svc.VerifyNewAuthUserOTP(context)

			testutil.RequireHttpResultJson(t, tt.expected, rec)
		})
	}
}
