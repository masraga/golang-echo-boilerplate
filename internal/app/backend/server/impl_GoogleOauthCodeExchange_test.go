package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/internal/app/backend/server"
)

func TestServer_GoogleOauthCodeExchange(t *testing.T) {
	type test struct {
		name         string
		expectedCode int
	}

	tests := []test{
		{
			name:         "not implemented",
			expectedCode: http.StatusNotImplemented,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			svc := server.NewServer(server.ServerOpts{})

			if err := svc.GoogleOauthCodeExchange(ctx); err != nil {
				t.Fatal(err)
			}
			if rec.Code != tt.expectedCode {
				t.Fatalf("expected status %d, got %d", tt.expectedCode, rec.Code)
			}
		})
	}
}
