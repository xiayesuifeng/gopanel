package server

import (
	"gitlab.com/xiayesuifeng/gopanel/api/server/endpoint/app"
	"gitlab.com/xiayesuifeng/gopanel/api/server/endpoint/auth"
	"gitlab.com/xiayesuifeng/gopanel/api/server/endpoint/backend"
	"gitlab.com/xiayesuifeng/gopanel/api/server/endpoint/caddy"
	"gitlab.com/xiayesuifeng/gopanel/api/server/endpoint/service"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
)

func (s *Server) registerAll() {
	endpoints := []router.Endpoint{
		&auth.Auth{},
		&app.App{},
		&backend.Backend{},
		&caddy.Caddy{},
		&service.Service{},
	}

	for _, endpoint := range endpoints {
		s.Register(endpoint)
	}
}
