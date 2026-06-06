package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/crypto"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/util/parser"
	"github.com/masraga/golang-echo-boilerplate/internal/util/pointer"
)

func (s *Server) AuthValidatePin(ctx echo.Context) error {
	var body api.AuthValidatePinRequest
	err := bindOrReturnBadRequest(ctx, &body)
	if err != nil {
		return err
	}

	phone, err := s.CryptoService.Decrypt(ctx.Request().Context(), crypto.DecryptInput{
		HashCode: body.PhoneNo,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	svc, err := s.AuthService.AuthValidatePin(ctx.Request().Context(), auth.AuthValidatePinInput{
		PhoneNo:       phone.Result,
		PinCode:       body.Pin,
		RetypePinCode: body.RetypePin,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	userId, err := parser.ParseToUUID(svc.UserId)
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, api.AuthValidatePinResponse{
		AuthToken: pointer.String(svc.Token),
		UserId:    userId,
	})
}
