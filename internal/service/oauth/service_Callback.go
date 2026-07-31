package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

func (s *OAuthService) Callback(ctx context.Context, input GoogleOauthCallbackInput) (output GoogleOauthCallbackOutput, err error) {
	providerOutput, err := s.provider.Callback(ctx, input)
	if err != nil {
		err = s.err.Wrap(errors.Join(err, ErrFailedGetCallbackToken))
		return
	}

	id := uuid.New().String()
	_, err = s.oauthRepositoryWriter.CreateAccessToken(ctx, CreateAccessTokenInput{
		Id:           id,
		AccessToken:  providerOutput.Token,
		RefreshToken: providerOutput.RefreshToken,
	})

	if err != nil {
		err = s.err.Wrap(err)
		return
	}

	fmt.Printf("%+v\n", providerOutput)

	// do something with callback state like register user, etc
	callbackState, err := s.convertStateToQueryString(input)
	if err != nil {
		err = s.err.Wrap(err)
		return
	}
	if len(callbackState.Actions) > 0 {
		for _, act := range callbackState.Actions {
			if act == ExchangeTokenActionRegisterUser {
				// since the callback from google is from valid user, so system must be
				// bypassing to insert new user data based on google token
				authUser, err := s.authService.FindAuthWithEmail(ctx, auth.FindAuthWithEmailInput{
					Email: providerOutput.Email,
				})
				// return error but no result set. and prevent add user
				// when user already exists
				if err != nil {
					if err == auth.ErrFindAuthWithEmail {
						return output, err
					}
				}
				if authUser.Email == providerOutput.Email {
					return providerOutput, nil
				}

				userId := uuid.NewString()
				s.oauthRepositoryWriter.BypassCreateNewUser(ctx, BypassCreateNewUserInput{
					Id:        userId,
					PhoneNo:   auth.DEFAULT_PHONE_NO,
					Pin:       auth.DEFAULT_PIN_CODE,
					Email:     providerOutput.Email, // since email is already provided from provider oauth, so we can pass the output
					CreatedBy: DefaultCreatedBy,
					LoginMode: auth.LOGIN_MODE_GOOGLE,
				})
			}
		}
	}

	output.Token = providerOutput.Token
	output.RefreshToken = providerOutput.RefreshToken
	output.State = providerOutput.State

	return
}

func (s *OAuthService) convertStateToQueryString(input GoogleOauthCallbackInput) (output GoogleOauthCallbackState, err error) {
	urlVal, err := url.ParseQuery(input.State)
	if err != nil {
		err = s.err.Wrap(errors.Join(err, ErrParseOauthAction))
		return
	}
	var state GoogleOauthCallbackState
	json.Unmarshal([]byte(urlVal.Get("actions")), &state.Actions)

	output.Actions = state.Actions
	return
}
