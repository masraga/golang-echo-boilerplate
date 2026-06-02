package auth

import (
	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/ctxerr"
)

type AuthService struct {
	JwtSecret                 JwtSecretType
	JwtExpiration             JwtExpirationType
	AuthAccessBootstrapUserId AuthAccessBootstrapUserIdType
	Err                       *ctxerr.CtxErr

	AuthRepositoryWriter AuthRepositoryWriterInterface
	AuthRepositoryReader AuthRepositoryReaderInterface
	CryptoService        crypto.CryptoServiceInterface
}

type AuthServiceOpts struct {
	JwtSecret                 JwtSecretType
	JwtExpiration             JwtExpirationType
	AuthAccessBootstrapUserId AuthAccessBootstrapUserIdType
	Err                       *ctxerr.CtxErr

	AuthRepositoryWriter AuthRepositoryWriterInterface
	AuthRepositoryReader AuthRepositoryReaderInterface
	CryptoService        crypto.CryptoServiceInterface
}

func NewAuthService(opts AuthServiceOpts) *AuthService {
	return &AuthService{
		JwtSecret:                 opts.JwtSecret,
		JwtExpiration:             opts.JwtExpiration,
		AuthAccessBootstrapUserId: opts.AuthAccessBootstrapUserId,
		Err:                       opts.Err,

		AuthRepositoryWriter: opts.AuthRepositoryWriter,
		AuthRepositoryReader: opts.AuthRepositoryReader,
		CryptoService:        opts.CryptoService,
	}
}
