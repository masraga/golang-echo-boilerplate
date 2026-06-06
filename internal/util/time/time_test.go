package time_test

import (
	"testing"
	stdtime "time"

	utiltime "github.com/masraga/golang-echo-boilerplate/internal/util/time"
	"github.com/stretchr/testify/require"
)

func TestNowUtc0(t *testing.T) {
	type test struct {
		name string
	}

	tests := []test{
		{
			name: "should return current unix millisecond time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := stdtime.Now().UnixMilli()
			got := utiltime.NowUtc0()
			after := stdtime.Now().UnixMilli()

			require.GreaterOrEqual(t, got, before)
			require.LessOrEqual(t, got, after)
		})
	}
}
