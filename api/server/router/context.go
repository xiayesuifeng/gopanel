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
