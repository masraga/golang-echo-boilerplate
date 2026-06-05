package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/service/auth"
)

func (s *Server) VerifyNewAuthUserOTP(ctx echo.Context) error {
	var reqBody api.VerifyOTPRequest
	err := bindOrReturnBadRequest(ctx, &reqBody)
	if err != nil {
		return err
	}

	phone, err := s.CryptoService.Decrypt(ctx.Request().Context(), crypto.DecryptInput{
		HashCode: reqBody.PhoneNo,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	output, err := s.AuthService.VerifyOtp(ctx.Request().Context(), auth.VerifyOtpInput{
		PhoneNo: phone.Result,
		OtpCode: reqBody.Otp,
	})
	if err != nil {
		return returnError(ctx, err)
	}

	var resp api.VerifyOTPResponse200
	resp.IsValid = output.IsValid
	resp.PhoneNo = output.PhoneNo
	resp.Note = output.Note
	resp.IsNewUser = output.IsNewUser

	return returnOk(ctx, resp)
}
