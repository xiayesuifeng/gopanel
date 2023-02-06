package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Context struct {
	*gin.Context
}

func (ctx *Context) JSON(data interface{}) error {
	ctx.Context.JSON(http.StatusOK, data)

	return nil
}

func (ctx *Context) NoContent() error {
	ctx.Context.Status(http.StatusNoContent)

	return nil
}

func (ctx *Context) Error(code int, message string) error {
	return &APIError{Message: message, Code: code}
}
