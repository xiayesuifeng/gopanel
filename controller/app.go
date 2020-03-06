package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/app"
)

type App struct {
}

func (a *App) Get(ctx *gin.Context) {

}

func (a *App) Gets(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"apps": app.GetApps(),
	})
}

func (a *App) Post(ctx *gin.Context) {
	data := app.App{}
	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	if err := app.AddApp(data); err != nil {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{"code": 200})
	}
}

func (a *App) Put(ctx *gin.Context) {

}

func (a *App) Delete(ctx *gin.Context) {

}
