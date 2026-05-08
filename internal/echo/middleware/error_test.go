package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHTTPErrorHandlerLogsStacktrace(t *testing.T) {
	e := echo.New()
	logs := new(strings.Builder)
	e.Logger.SetOutput(logs)

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/error")
	c.Response().Header().Set(echo.HeaderXRequestID, "request-id")

	HTTPErrorHandler(errors.New("service failed"), c)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	log := logs.String()
	for _, want := range []string{
		"request error:",
		"method=GET",
		"path=/error",
		"route=/error",
		"request_id=request-id",
		"error=service failed",
		"origin=unknown",
		"stacktrace:",
		"runtime/debug.Stack",
	} {
		if !strings.Contains(log, want) {
			t.Fatalf("expected log to contain %q, got %s", want, log)
		}
	}
}
