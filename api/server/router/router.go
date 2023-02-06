package router

import "github.com/gin-gonic/gin"

type Route interface{}

type HandlerFunc func(ctx *Context) error

type Endpoint interface {
	Name() string
	Run(r Router)
}

type Router interface {
	Routes
	Group(string, ...HandlerFunc) Router
}

type Routes interface {
	Handle(string, string, ...HandlerFunc)
	Any(string, ...HandlerFunc)
	GET(string, ...HandlerFunc)
	POST(string, ...HandlerFunc)
	DELETE(string, ...HandlerFunc)
	PATCH(string, ...HandlerFunc)
	PUT(string, ...HandlerFunc)
	OPTIONS(string, ...HandlerFunc)
	HEAD(string, ...HandlerFunc)
}

func NewRouter(router gin.IRouter) Router {
	return &wrapperRouter{router}
}
