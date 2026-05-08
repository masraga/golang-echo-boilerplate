package testutil

import (
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
