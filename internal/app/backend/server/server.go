package server

import (
	"github.com/masraga/kerp-api/internal/service/auth"
)

type Server struct {
	AuthService auth.AuthServiceInterface
}

func NewServer(
	authService auth.AuthServiceInterface,
) *Server {
	return &Server{
		AuthService: authService,
	}
}
