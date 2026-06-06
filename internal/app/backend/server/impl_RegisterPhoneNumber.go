package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/crypto"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

func (s *Server) RegisterPhoneNumber(ctx echo.Context) error {
	input, err := s.bindRequestToCreateNewAccountInput(ctx)
	if err != nil {
		return returnBadRequest(ctx, err.Error())
	}
	svc, err := s.AuthService.CreateNewAccount(ctx.Request().Context(), input)
	if err != nil {
		return returnError(ctx, err)
	}

	output, err := s.bindOutputWithCreateNewAccountOutput(svc)

	return returnCreated(ctx, output)
}

func (s *Server) bindRequestToCreateNewAccountInput(ctx echo.Context) (input auth.CreateNewAccountInput, err error) {
	var reqBody api.CreateNewAccountRequest
	err = bindOrReturnBadRequest(ctx, &reqBody)
	if err != nil {
		return
	}
	phone, err := s.CryptoService.Decrypt(ctx.Request().Context(), crypto.DecryptInput{
		HashCode: reqBody.PhoneNo,
	})
	if err != nil {
		return
	}
	input.PhoneNo = phone.Result
	input.FirebaseId = &reqBody.FirebaseId
	return
}

func (s *Server) bindOutputWithCreateNewAccountOutput(result auth.CreateNewAccountOutput) (output api.CreateNewAccountResponse, err error) {
	output.Id = result.Id
	output.OtpCode = result.OtpCode
	return
}
