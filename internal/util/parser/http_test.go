package parser_test

import (
	"testing"

	"github.com/masraga/kerp-api/internal/util/parser"
	"github.com/stretchr/testify/require"
)

func TestNormalizeEndpointPath(t *testing.T) {
	type args struct {
		path string
	}

	type test struct {
		name     string
		args     args
		expected string
	}

	tests := []test{
		{
			name:     "should preserve empty path",
			args:     args{path: ""},
			expected: "",
		},
		{
			name:     "should trim and prefix path",
			args:     args{path: " api/v1/auth/roles "},
			expected: "/api/v1/auth/roles",
		},
		{
			name:     "should preserve prefixed path",
			args:     args{path: "/api/v1/auth/roles"},
			expected: "/api/v1/auth/roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.NormalizeEndpointPath(tt.args.path)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestNormalizeEndpointMethod(t *testing.T) {
	type args struct {
		method string
	}

	type test struct {
		name     string
		args     args
		expected string
	}

	tests := []test{
		{
			name:     "should trim and lowercase method",
			args:     args{method: " POST "},
			expected: "post",
		},
		{
			name:     "should preserve empty method",
			args:     args{method: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.NormalizeEndpointMethod(tt.args.method)
			require.Equal(t, tt.expected, got)
		})
	}
}
