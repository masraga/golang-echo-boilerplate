package auth

import (
	"github.com/masraga/kerp-api/internal/ctxerr"
)

type AuthService struct {
	JwtSecret     JwtSecretType
	JwtExpiration JwtExpirationType
	Err           *ctxerr.CtxErr

	AuthRepositoryWriter AuthRepositoryWriterInterface
	AuthRepositoryReader AuthRepositoryReaderInterface
}

type AuthServiceOpts struct {
	JwtSecret     JwtSecretType
	JwtExpiration JwtExpirationType
	Err           *ctxerr.CtxErr

	AuthRepositoryWriter AuthRepositoryWriterInterface
	AuthRepositoryReader AuthRepositoryReaderInterface
}

func NewAuthService(opts AuthServiceOpts) *AuthService {
	return &AuthService{
		JwtSecret:     opts.JwtSecret,
		JwtExpiration: opts.JwtExpiration,
		Err:           opts.Err,

		AuthRepositoryWriter: opts.AuthRepositoryWriter,
		AuthRepositoryReader: opts.AuthRepositoryReader,
	}
}
