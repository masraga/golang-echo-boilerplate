package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/crypto"
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
