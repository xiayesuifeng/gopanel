package server

import (
	"gitlab.com/xiayesuifeng/gopanel/api/server/endpoint/auth"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
)

func (s *Server) registerAll() {
	endpoints := []router.Endpoint{
		&auth.Auth{},
	}

	for _, endpoint := range endpoints {
		s.Register(endpoint)
	}
}
