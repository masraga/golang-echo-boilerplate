package oauth

import "errors"

var (
	ErrFailedExchangeToken    error = errors.New("failed to exchange authorization code")
	ErrFailedGetCallbackToken error = errors.New("failed to get token callback")
	ErrAccessTokenNotFound    error = errors.New("cant find access token")
	ErrCreateNewAccessToken   error = errors.New("error occur when create oauth access token")
	ErrParseOauthAction       error = errors.New("error when parse state action to url")
	ErrBypassCreateUser       error = errors.New("error to bypass new user")
)
