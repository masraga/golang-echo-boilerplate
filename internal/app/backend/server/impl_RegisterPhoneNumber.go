package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/service/auth"
)

func (s *Server) RegisterPhoneNumber(ctx echo.Context) error {
	input, err := s.bindRequestToCreateNewAccountInput(ctx)
	if err != nil {
		return s.returnBadRequest(ctx, err.Error())
	}
	svc, err := s.AuthService.CreateNewAccount(ctx.Request().Context(), input)
	if err != nil {
		return s.returnError(ctx, err)
	}

	output, err := s.bindOutputWithCreateNewAccountOutput(svc)

	return s.returnCreated(ctx, output)
}

func (s *Server) bindRequestToCreateNewAccountInput(ctx echo.Context) (input auth.CreateNewAccountInput, err error) {
	var reqBody api.CreateNewAccountRequest
	err = s.bindOrReturnBadRequest(ctx, &reqBody)
	if err != nil {
		return
	}
	input.PhoneNo = reqBody.PhoneNo
	return
}

func (s *Server) bindOutputWithCreateNewAccountOutput(result auth.CreateNewAccountOutput) (output api.CreateNewAccountResponse, err error) {
	output.Id = result.Id
	return
}
