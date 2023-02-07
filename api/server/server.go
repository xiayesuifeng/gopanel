package server

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/middleware"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
)

type Server struct {
	endpoints []router.Endpoint
}

func NewServer() *Server {
	server := &Server{}

	server.registerAll()

	return server
}

func (s *Server) Register(endpoint router.Endpoint) {
	s.endpoints = append(s.endpoints, endpoint)
}

func (s *Server) Run(address ...string) error {
	engine := gin.Default()

	apiRouter := engine.Group("/api")

	r := router.NewRouter(apiRouter)
	r.Use(middleware.AuthMiddleware)

	for _, e := range s.endpoints {
		e.Run(r.Group("/" + e.Name()))
	}

	return engine.Run(address...)
}
