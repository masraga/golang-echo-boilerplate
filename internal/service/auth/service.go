package auth

import "github.com/masraga/kerp-api/internal/dbtx"

type AuthService struct {
	JwtSecret     JwtSecretType
	JwtExpiration JwtExpirationType

	dbtx.DbTxInterface
	AuthRepositoryWriter AuthRepositoryWriterInterface
	AuthRepositoryReader AuthRepositoryReaderInterface
}

type AuthServiceOpts struct {
	JwtSecret     JwtSecretType
	JwtExpiration JwtExpirationType

	dbtx.DbTxInterface
	AuthRepositoryWriter AuthRepositoryWriterInterface
	AuthRepositoryReader AuthRepositoryReaderInterface
}

func NewAuthService(opts AuthServiceOpts) *AuthService {
	return &AuthService{
		JwtSecret:            opts.JwtSecret,
		JwtExpiration:        opts.JwtExpiration,
		DbTxInterface:        opts.DbTxInterface,
		AuthRepositoryWriter: opts.AuthRepositoryWriter,
		AuthRepositoryReader: opts.AuthRepositoryReader,
	}
}
