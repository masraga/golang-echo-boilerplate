package oauth

type OAuthService struct {
	provider OAuthProviderInterface

	// oauthRepositoryReader OAuthRepositoryReaderInterface
	// oauthRepositoryWriter OAuthRepositoryWriterInterface
}

type OAuthServiceOpts struct {
	Provider OAuthProviderInterface

	// OauthRepositoryReader OAuthRepositoryReaderInterface
	// OauthRepositoryWriter OAuthRepositoryWriterInterface
}

func NewOAuthService(opts OAuthServiceOpts) *OAuthService {
	return &OAuthService{
		provider: opts.Provider,

		// oauthRepositoryReader: opts.OauthRepositoryReader,
		// oauthRepositoryWriter: opts.OauthRepositoryWriter,
	}
}
