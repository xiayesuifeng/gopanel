package firewall

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/firewall"
)

type Firewall struct {
}

func (f *Firewall) Name() string {
	return "firewall"
}

func (f *Firewall) Run(r router.Router) {
	r.Use(permanentMiddleware)

	r.GET("", f.GetConfig)
}

func (f *Firewall) GetConfig(ctx *router.Context) error {
	zone, err := firewall.GetDefaultZone()
	if err != nil {
		return err
	}

	return ctx.JSON(gin.H{
		"defaultZone": zone,
	})
}
