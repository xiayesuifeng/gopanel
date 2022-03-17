package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/configuration/caddy"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
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

func (c *Caddy) PutConfiguration(ctx *gin.Context) {
	configuration := caddy.Configuration{}

	if err := ctx.ShouldBind(&configuration); err != nil {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})

		return
	}

	if err := caddy.SetConfiguration(&configuration); err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	} else {
		caddyManager.GetManager().NotifyCaddyConfigChange()

		ctx.JSON(200, gin.H{"code": 200})
	}
}
