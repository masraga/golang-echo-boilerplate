package server

import (
	"github.com/masraga/golang-echo-boilerplate/internal/crypto"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
	"github.com/masraga/golang-echo-boilerplate/internal/service/notification"
	"github.com/masraga/golang-echo-boilerplate/internal/service/oauth"
)

type Server struct {
	AuthService         auth.AuthServiceInterface
	CryptoService       crypto.CryptoServiceInterface
	NotificationService notification.NotificationServiceInterface
	OAuthService        oauth.OAuthServiceInterface
}

type ServerOpts struct {
	AuthService         auth.AuthServiceInterface
	CryptoService       crypto.CryptoServiceInterface
	NotificationService notification.NotificationServiceInterface
	OAuthService        oauth.OAuthServiceInterface
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		AuthService:         opts.AuthService,
		CryptoService:       opts.CryptoService,
		NotificationService: opts.NotificationService,
		OAuthService:        opts.OAuthService,
	}
}
