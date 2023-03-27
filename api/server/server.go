package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/middleware"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"golang.org/x/exp/slog"
	"io/fs"
	"net/http"
	"time"
)

type Server struct {
	server    http.Server
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

func (s *Server) Run(address string) error {
	engine := gin.Default()

	s.registerWeb(engine)

	apiRouter := engine.Group("/api")

	r := router.NewRouter(apiRouter)
	r.Use(middleware.InstallMiddleware)
	r.Use(middleware.AuthMiddleware)

	for _, e := range s.endpoints {
		e.Run(r.Group("/" + e.Name()))
	}

	s.server = http.Server{
		Addr:    address,
		Handler: engine,
	}

	slog.Info("[server] listening and serving HTTP on " + address)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("[server] waiting for server shutdown...")
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("[server] failed to shutdown", err)
		return err
	}

	slog.Info("[server] server shutdown")

	return nil
}
