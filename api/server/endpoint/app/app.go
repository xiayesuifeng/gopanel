package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"gitlab.com/xiayesuifeng/gopanel/module/app"
	"strings"
)

type App struct {
}

func (a *App) Name() string {
	return "app"
}

func (a *App) Run(r router.Router) {
	r.GET("", a.Gets)
	r.GET("/:name", a.Get)
	r.POST("", a.Post)
	r.PUT("/:name", a.Put)
	r.DELETE("/:name", a.Delete)
}

func (a *App) GetValidate(ctx *router.Context) bool {
	validate := false
	host := config.Conf.Panel.Domain
	if host == "" {
		if strings.HasSuffix(ctx.Request.Host, fmt.Sprintf(":%d", config.Conf.Panel.Port)) {
			validate = true
		}
	} else if config.Conf.Panel.Port == 0 {
		if strings.HasPrefix(ctx.Request.Host, host) {
			validate = true
		}
	} else {
		host += fmt.Sprintf(":%d", config.Conf.Panel.Port)
		validate = ctx.Request.Host == host
	}

	return validate
}

func (a *App) Get(ctx *router.Context) error {
	name := ctx.Param("name")

	if app, err := app.GetApp(name); err == nil {
		return ctx.JSON(gin.H{
			"app": app,
		})
	} else if err.Error() == "app not found" {
		return ctx.Error(404, "app not found")
	} else {
		return err
	}
}

func (a *App) Gets(ctx *router.Context) error {
	return ctx.JSON(gin.H{
		"apps": app.GetApps(),
	})
}

func (a *App) Post(ctx *router.Context) error {
	data := app.App{}
	if err := ctx.ShouldBind(&data); err != nil {
		return ctx.Error(400, err.Error())
	}

	if data.Name == "" {
		return ctx.Error(400, "name must exist")
	}

	if err := app.AddApp(data, a.GetValidate(ctx)); err != nil {
		return ctx.Error(400, err.Error())
	} else {
		return ctx.NoContent()
	}
}

func (a *App) Put(ctx *router.Context) error {
	name := ctx.Param("name")

	data := app.App{}
	if err := ctx.ShouldBind(&data); err != nil {
		return ctx.Error(400, err.Error())
	}

	data.Name = name

	if err := app.EditApp(data, a.GetValidate(ctx)); err != nil {
		return err
	} else {
		return ctx.NoContent()
	}
}

func (a *App) Delete(ctx *router.Context) error {
	name := ctx.Param("name")

	if err := app.DeleteApp(name); err == nil {
		return ctx.NoContent()
	} else if err.Error() == "app not found" {
		return ctx.Error(404, "app not found")
	} else {
		return err
	}
}
