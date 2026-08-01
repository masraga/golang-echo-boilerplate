package oauth

import (
	"github.com/masraga/golang-echo-boilerplate/internal/ctxerr"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/rs/zerolog"
)

type OAuthService struct {
	provider    OAuthProviderInterface
	authService auth.AuthServiceInterface

	logger                zerolog.Logger
	err                   *ctxerr.CtxErr
	oauthRepositoryReader OAuthRepositoryReaderInterface
	oauthRepositoryWriter OAuthRepositoryWriterInterface
}

type OAuthServiceOpts struct {
	Provider    OAuthProviderInterface
	AuthService auth.AuthServiceInterface

	Logger                zerolog.Logger
	Err                   *ctxerr.CtxErr
	OauthRepositoryReader OAuthRepositoryReaderInterface
	OauthRepositoryWriter OAuthRepositoryWriterInterface
}

func NewOAuthService(opts OAuthServiceOpts) *OAuthService {
	return &OAuthService{
		provider:    opts.Provider,
		authService: opts.AuthService,

		logger:                opts.Logger,
		err:                   opts.Err,
		oauthRepositoryReader: opts.OauthRepositoryReader,
		oauthRepositoryWriter: opts.OauthRepositoryWriter,
	}
}
