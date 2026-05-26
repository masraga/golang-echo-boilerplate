package server

import (
	"github.com/masraga/kerp-api/internal/crypto"
	"github.com/masraga/kerp-api/internal/service/auth"
)

type Server struct {
	AuthService   auth.AuthServiceInterface
	CryptoService crypto.CryptoServiceInterface
}

type ServerOpts struct {
	AuthService   auth.AuthServiceInterface
	CryptoService crypto.CryptoServiceInterface
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		AuthService:   opts.AuthService,
		CryptoService: opts.CryptoService,
	}
}
