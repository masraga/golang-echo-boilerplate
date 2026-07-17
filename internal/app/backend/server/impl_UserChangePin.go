package server

import (
	"github.com/labstack/echo/v4"
	api "github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/crypto"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

func (s *Server) UserChangePin(ctx echo.Context) error {
	var reqBody api.UserChangePinRequest
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	userPhoneNo, err := s.CryptoService.Decrypt(ctx.Request().Context(), crypto.DecryptInput{
		HashCode: reqBody.UserPhoneNo,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	oldPin, err := s.CryptoService.Decrypt(ctx.Request().Context(), crypto.DecryptInput{
		HashCode: reqBody.OldPin,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	newPin, err := s.CryptoService.Decrypt(ctx.Request().Context(), crypto.DecryptInput{
		HashCode: reqBody.NewPin,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	retypeNewPin, err := s.CryptoService.Decrypt(ctx.Request().Context(), crypto.DecryptInput{
		HashCode: reqBody.RetypeNewPin,
	})
	if err != nil {
		return returnError(ctx, err)
	}

	svc, err := s.AuthService.UserChangePin(ctx.Request().Context(), auth.UserChangePinInput{
		UserPhoneNo:  userPhoneNo.Result,
		OldPin:       oldPin.Result,
		NewPin:       newPin.Result,
		RetypeNewPin: retypeNewPin.Result,
	})
	if err != nil {
		return returnError(ctx, err)
	}

	return returnOk(ctx, api.UserChangePinResponse{
		IsUpdate: svc.IsUpdate,
	})
}
