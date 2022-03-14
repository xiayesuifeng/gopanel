package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/configuration/caddy"
)

type Caddy struct {
}

func (c *Caddy) GetConfiguration(ctx *gin.Context) {
	conf := caddy.GetConfiguration()
	ctx.JSON(200, gin.H{
		"code":          200,
		"configuration": conf,
	})
}
