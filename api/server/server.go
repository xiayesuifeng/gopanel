package server

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/middleware"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"io/fs"
)

type Server struct {
	endpoints []router.Endpoint
	web       fs.FS
}

func NewServer(web fs.FS) *Server {
	server := &Server{web: web}

	server.registerAll()

	return server
}

func (s *Server) Register(endpoint router.Endpoint) {
	s.endpoints = append(s.endpoints, endpoint)
}

func (s *Server) Run(address ...string) error {
	engine := gin.Default()

	s.registerWeb(engine)

	apiRouter := engine.Group("/api")

	r := router.NewRouter(apiRouter)
	r.Use(middleware.AuthMiddleware)

	for _, e := range s.endpoints {
		e.Run(r.Group("/" + e.Name()))
	}

	return engine.Run(address...)
}
