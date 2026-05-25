package testutil

import (
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Result[v any] struct {
	Err   error
	Value v
}

func RequireResult[v any](t require.TestingT, err error, expected Result[v], actual v, msgAndArgs ...interface{}) {
	if expected.Err != nil {
		require.ErrorIs(t, err, expected.Err, msgAndArgs...)
		return
	}
	require.NoError(t, err, msgAndArgs...)
	require.Equal(t, expected.Value, actual, msgAndArgs...)
}

type HttpResult struct {
	Code int
	Body string
}

func RequireHttpResultJson(t require.TestingT, expected HttpResult, actual *httptest.ResponseRecorder) {
	assert.Equal(t, expected.Code, actual.Code)
	assert.JSONEq(t, expected.Body, actual.Body.String())
}
