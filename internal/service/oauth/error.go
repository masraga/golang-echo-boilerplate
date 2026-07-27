package oauth

import "errors"

var (
	ErrFailedExchangeToken    error = errors.New("failed to exchange authorization code")
	ErrFailedGetCallbackToken error = errors.New("failed to get token callback")
	ErrAccessTokenNotFound    error = errors.New("cant find access token")
)
