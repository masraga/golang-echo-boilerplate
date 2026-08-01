package gauth

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
)

func (s *GAuthService) Callback(ctx context.Context, input oauth.GoogleOauthCallbackInput) (output oauth.GoogleOauthCallbackOutput, err error) {
	oauthConfig, err := GenerateGauthConfig(GenerateGauthConfigInput{
		ClientId:        s.clientId,
		ClientSecret:    s.clientSecret,
		AuthRedirectUrl: s.authRedirectUrl,
		Scopes:          &s.scopes,
	})
	if err != nil {
		err = s.err.Wrap(err)
		return
	}

	token, err := oauthConfig.Exchange(ctx, input.Code)
	if err != nil {
		return
	}

	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		err = s.err.Wrap(errors.Join(err, oauth.ErrFailedGetCallbackToken))
		return
	}

	oidcProvider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		err = s.err.Wrap(errors.Join(err, oauth.ErrFailedGetCallbackToken))
		return
	}

	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID: string(s.clientId),
	})
	idToken, err := verifier.Verify(ctx, idTokenStr)
	if err != nil {
		err = s.err.Wrap(errors.Join(err, oauth.ErrFailedGetCallbackToken))
		return
	}
	if err = idToken.Claims(&output); err != nil {
		err = s.err.Wrap(errors.Join(err, oauth.ErrFailedGetCallbackToken))
		return
	}

	output.IdToken = idTokenStr
	output.Token = token.AccessToken
	output.RefreshToken = token.RefreshToken
	output.State = input.State

	return
}
