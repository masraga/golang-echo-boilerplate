package parser

import "strings"

func NormalizeEndpointPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" || strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func NormalizeEndpointMethod(method string) string {
	return strings.ToLower(strings.TrimSpace(method))
}
