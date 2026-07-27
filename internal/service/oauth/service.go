package oauth

import "github.com/masraga/golang-echo-boilerplate/internal/ctxerr"

type OAuthService struct {
	provider OAuthProviderInterface

	err                   *ctxerr.CtxErr
	oauthRepositoryReader OAuthRepositoryReaderInterface
	oauthRepositoryWriter OAuthRepositoryWriterInterface
}

type OAuthServiceOpts struct {
	Provider OAuthProviderInterface

	Err                   *ctxerr.CtxErr
	OauthRepositoryReader OAuthRepositoryReaderInterface
	OauthRepositoryWriter OAuthRepositoryWriterInterface
}

func NewOAuthService(opts OAuthServiceOpts) *OAuthService {
	return &OAuthService{
		provider: opts.Provider,

		err:                   opts.Err,
		oauthRepositoryReader: opts.OauthRepositoryReader,
		oauthRepositoryWriter: opts.OauthRepositoryWriter,
	}
}
