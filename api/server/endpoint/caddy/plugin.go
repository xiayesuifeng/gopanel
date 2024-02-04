package caddy

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyutil/caddymodule"
)

func (c *Caddy) GetModuleList(ctx *router.Context) error {
	if list, err := caddymodule.GetModuleList(); err != nil {
		return err
	} else {
		return ctx.JSON(gin.H{
			"modules": list,
		})
	}
}

func (c *Caddy) GetOfficialPluginList(ctx *router.Context) error {
	list, err := caddymodule.GetOfficialPluginList()
	if err != nil {
		return err
	}

	return ctx.JSON(list)
}

func (c *Caddy) InstallPlugin(ctx *router.Context) error {
	data := &struct {
		Packages []string `json:"packages"`
	}{}

	if err := ctx.ShouldBind(data); err != nil {
		return err
	}

	if err := caddymodule.InstallPlugin(data.Packages...); err != nil {
		return err
	}

	return ctx.NoContent()
}

func (c *Caddy) RemovePlugin(ctx *router.Context) error {
	pkgs := ctx.QueryArray("package")

	if err := caddymodule.RemovePlugin(pkgs...); err != nil {
		return err
	}

	return ctx.NoContent()
}
