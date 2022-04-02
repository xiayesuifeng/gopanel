package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyutil/caddyconfig"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyutil/caddymodule"
)

type Caddy struct {
}

func (c *Caddy) GetConfiguration(ctx *gin.Context) {
	conf := caddyconfig.GetConfiguration()
	ctx.JSON(200, gin.H{
		"code":          200,
		"configuration": conf,
	})
}

func (c *Caddy) PutConfiguration(ctx *gin.Context) {
	configuration := caddyconfig.Configuration{}

	if err := ctx.ShouldBind(&configuration); err != nil {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})

		return
	}

	if err := caddyconfig.SetConfiguration(&configuration); err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	} else {
		caddyManager.GetManager().NotifyCaddyConfigChange()

		ctx.JSON(200, gin.H{"code": 200})
	}
}

func (c *Caddy) GetModuleList(ctx *gin.Context) {
	if list, err := caddymodule.GetModuleList(); err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"modules": list,
		})
	}
}
