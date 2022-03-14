package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/configuration/caddy"
)

type CaddyConf struct {
}

func (c *CaddyConf) Get(ctx *gin.Context) {
	conf := caddy.GetConfiguration()
	ctx.JSON(200, gin.H{
		"code":          200,
		"configuration": conf,
	})
}
