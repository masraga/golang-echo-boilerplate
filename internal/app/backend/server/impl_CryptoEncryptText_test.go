package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/app/backend/server"
	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestServer_CryptoEncryptText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	requestBody := api.CryptoEncryptTextRequest{
		PlainText: "081234567890",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto/encrypt", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	cryptoService := crypto.NewMockCryptoServiceInterface(ctrl)
	cryptoService.EXPECT().
		Encrypt(ctx.Request().Context(), crypto.EncryptInput{PlainText: requestBody.PlainText}).
		Return(crypto.EncryptOutput{Result: "encrypted-code"}, nil)

	svc := server.NewServer(server.ServerOpts{
		CryptoService: cryptoService,
	})
	err = svc.CryptoEncryptText(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := json.Marshal(api.CryptoEncryptTextResponse{Code: "encrypted-code"})
	if err != nil {
		t.Fatal(err)
	}
	testutil.RequireHttpResultJson(t, testutil.HttpResult{
		Code: http.StatusOK,
		Body: string(expected),
	}, rec)
}
