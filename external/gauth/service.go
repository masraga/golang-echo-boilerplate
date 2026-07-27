package gauth

import "github.com/masraga/golang-echo-boilerplate/internal/ctxerr"

type GAuthService struct {
	clientSecret    ClientSecret
	clientId        ClientId
	authRedirectUrl AuthRedirectUrl
	scopes          Scopes
	err             *ctxerr.CtxErr
}

type GAuthServiceOpts struct {
	ClientSecret    ClientSecret
	ClientId        ClientId
	AuthRedirectUrl AuthRedirectUrl
	Scopes          Scopes
	Err             *ctxerr.CtxErr
}

func NewGAuthService(opts GAuthServiceOpts) *GAuthService {
	return &GAuthService{
		clientSecret:    opts.ClientSecret,
		clientId:        opts.ClientId,
		authRedirectUrl: opts.AuthRedirectUrl,
		scopes:          opts.Scopes,
		err:             opts.Err,
	}
}
