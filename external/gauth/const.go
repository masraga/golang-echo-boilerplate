package gauth

import (
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GenerateGauthConfig(input GenerateGauthConfigInput) (output *oauth2.Config, err error) {
	var scopes []string
	if input.Scopes != nil {
		err := json.Unmarshal([]byte(*input.Scopes), &scopes)
		if err != nil {
			return nil, err
		}
	}
	config := &oauth2.Config{
		ClientID:     string(input.ClientId),
		ClientSecret: string(input.ClientSecret),
		RedirectURL:  string(input.AuthRedirectUrl),
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}

	return config, nil
}
