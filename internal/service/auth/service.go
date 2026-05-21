package auth

type AuthService struct {
	JwtSecret     JwtSecretType
	JwtExpiration JwtExpirationType

	AuthRepositoryWriter AuthRepositoryWriterInterface
	AuthRepositoryReader AuthRepositoryReaderInterface
}

type AuthServiceOpts struct {
	JwtSecret     JwtSecretType
	JwtExpiration JwtExpirationType

	AuthRepositoryWriter AuthRepositoryWriterInterface
	AuthRepositoryReader AuthRepositoryReaderInterface
}

func NewAuthService(opts AuthServiceOpts) *AuthService {
	return &AuthService{
		JwtSecret:            opts.JwtSecret,
		JwtExpiration:        opts.JwtExpiration,
		AuthRepositoryWriter: opts.AuthRepositoryWriter,
		AuthRepositoryReader: opts.AuthRepositoryReader,
	}
}
