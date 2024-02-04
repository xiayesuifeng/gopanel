package caddy

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyapp/caddyddns"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyutil/caddyconfig"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyutil/caddymodule"
)

type Caddy struct {
}

func (c *Caddy) Name() string {
	return "caddy"
}

func (c *Caddy) Run(r router.Router) {
	r.GET("/configuration", c.GetConfiguration)
	r.PUT("/configuration", c.PutConfiguration)

	r.GET("/module", c.GetModuleList)
	r.GET("/plugin/repo", c.GetOfficialPluginList)
	r.POST("/plugin", c.InstallPlugin)
	r.DELETE("/plugin", c.RemovePlugin)

	r.GET("/ddns", c.GetDynamicDNS)
	r.PUT("/ddns", c.PutDynamicDNS)
}

func (c *Caddy) GetConfiguration(ctx *router.Context) error {
	conf := caddyconfig.GetConfiguration()
	return ctx.JSON(gin.H{
		"configuration": conf,
	})
}

func (c *Caddy) PutConfiguration(ctx *router.Context) error {
	configuration := caddyconfig.Configuration{}

	if err := ctx.ShouldBind(&configuration); err != nil {
		return ctx.Error(400, err.Error())
	}

	if err := caddyconfig.SetConfiguration(&configuration); err != nil {
		return err
	}

	caddyManager.GetManager().NotifyCaddyConfigChange()

	return ctx.NoContent()
}

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

func (c *Caddy) GetDynamicDNS(ctx *router.Context) error {
	if caddyDDNS, err := caddyddns.GetCaddyDDNS(); err != nil {
		return err
	} else {
		return ctx.JSON(gin.H{
			"config": caddyDDNS,
		})
	}
}

func (c *Caddy) PutDynamicDNS(ctx *router.Context) error {
	caddyDDNS := caddyddns.CaddyDDNS{}

	if err := ctx.ShouldBind(&caddyDDNS); err != nil {
		return ctx.Error(400, err.Error())
	}

	if err := caddyddns.SetCaddyDDNS(&caddyDDNS); err != nil {
		return err
	} else {
		return ctx.NoContent()
	}
}
