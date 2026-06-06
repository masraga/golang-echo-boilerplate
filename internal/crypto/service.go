package crypto

import "github.com/masraga/golang-echo-boilerplate/internal/ctxerr"

type CryptoService struct {
	ConfigCryptoKey
	Err *ctxerr.CtxErr
}

type CryptoServiceOpts struct {
	ConfigCryptoKey
	Err *ctxerr.CtxErr
}

func NewCryptoService(opts CryptoServiceOpts) *CryptoService {
	return &CryptoService{
		ConfigCryptoKey: opts.ConfigCryptoKey,
		Err:             opts.Err,
	}
}
