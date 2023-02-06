package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type wrapperRouter struct {
	gin.IRouter
}

func (w *wrapperRouter) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) {
	w.IRouter.Handle(httpMethod, relativePath, w.wrapHandler(handlers...)...)
}

func (w *wrapperRouter) Any(relativePath string, handlers ...HandlerFunc) {
	w.IRouter.Any(relativePath, w.wrapHandler(handlers...)...)
}

func (w *wrapperRouter) GET(relativePath string, handlers ...HandlerFunc) {
	w.IRouter.GET(relativePath, w.wrapHandler(handlers...)...)
}

func (w *wrapperRouter) POST(relativePath string, handlers ...HandlerFunc) {
	w.IRouter.POST(relativePath, w.wrapHandler(handlers...)...)
}

func (w *wrapperRouter) DELETE(relativePath string, handlers ...HandlerFunc) {
	w.IRouter.DELETE(relativePath, w.wrapHandler(handlers...)...)
}

func (w *wrapperRouter) PATCH(relativePath string, handlers ...HandlerFunc) {
	w.IRouter.PATCH(relativePath, w.wrapHandler(handlers...)...)
}

func (w *wrapperRouter) PUT(relativePath string, handlers ...HandlerFunc) {
	w.IRouter.PUT(relativePath, w.wrapHandler(handlers...)...)
}

func (w *wrapperRouter) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	w.IRouter.OPTIONS(relativePath, w.wrapHandler(handlers...)...)
}

func (w *wrapperRouter) HEAD(relativePath string, handlers ...HandlerFunc) {
	w.IRouter.HEAD(relativePath, w.wrapHandler(handlers...)...)
}

func (w *wrapperRouter) Group(relativePath string, handlers ...HandlerFunc) Router {
	return &wrapperRouter{w.IRouter.Group(relativePath, w.wrapHandler(handlers...)...)}
}

func (w *wrapperRouter) wrapHandler(handlerFunc ...HandlerFunc) (funks []gin.HandlerFunc) {
	for _, f := range handlerFunc {
		funks = append(funks, func(ctx *gin.Context) {
			err := f(&Context{ctx})
			if err != nil {
				if apiError, ok := err.(*APIError); ok {
					ctx.JSON(apiError.Code, apiError)
				} else {
					ctx.JSON(http.StatusInternalServerError, APIError{Message: err.Error()})
				}
			}
		})
	}

	return
}
