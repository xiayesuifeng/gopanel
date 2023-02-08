package server

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

type embedFileSystem struct {
	http.FileSystem
	indexes bool
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	f, err := e.Open(path)
	if err != nil {
		return false
	}

	s, _ := f.Stat()
	if s.IsDir() && !e.indexes {
		return false
	}

	return true
}

func embedFile(webFS fs.FS, index bool) static.ServeFileSystem {
	return embedFileSystem{
		FileSystem: http.FS(webFS),
		indexes:    index,
	}
}

func (s *Server) registerWeb(engine *gin.Engine) {
	index := s.getWebIndex()

	engine.Use(static.Serve("/", embedFile(s.web, false)))
	engine.NoRoute(func(ctx *gin.Context) {
		if !strings.Contains(ctx.Request.RequestURI, "/api") && !strings.Contains(ctx.Request.RequestURI, "/netdata") {
			path := strings.Split(ctx.Request.URL.Path, "/")
			if len(path) > 1 && len(index) > 0 {
				ctx.Data(http.StatusOK, http.DetectContentType(index), index)
			}
		}
	})
}

func (s *Server) getWebIndex() []byte {
	file, err := s.web.Open("index.html")
	if err != nil {
		return []byte{}
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return []byte{}
	}

	return bytes
}
