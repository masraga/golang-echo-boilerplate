package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/crypto"
)

func (s *Server) CryptoEncryptText(ctx echo.Context) error {
	var reqBody api.CryptoEncryptTextRequest

	err := bindOrReturnBadRequest(ctx, &reqBody)
	if err != nil {
		return err
	}

	output, err := s.CryptoService.Encrypt(ctx.Request().Context(), crypto.EncryptInput{
		PlainText: reqBody.PlainText,
	})
	if err != nil {
		return returnError(ctx, err)
	}

	var result api.CryptoEncryptTextResponse
	result.Code = output.Result
	return returnOk(ctx, result)
}
