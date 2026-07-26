package gauth

type GAuthService struct {
	clientSecret    ClientSecret
	clientId        ClientId
	authRedirectUrl AuthRedirectUrl
}

type GAuthServiceOpts struct {
	ClientSecret    ClientSecret
	ClientId        ClientId
	AuthRedirectUrl AuthRedirectUrl
}

func NewGAuthService(opts GAuthServiceOpts) *GAuthService {
	return &GAuthService{
		clientSecret:    opts.ClientSecret,
		clientId:        opts.ClientId,
		authRedirectUrl: opts.AuthRedirectUrl,
	}
}
