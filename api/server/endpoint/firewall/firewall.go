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
	r.GET("/zone", f.GetZone)
	r.GET("/zone/names", f.GetZoneNames)
	r.GET("/zone/:name", f.GetZoneByName)
	r.PUT("/zone/:name", f.UpdateZoneByName)

	r.GET("/zone/:name/trafficRule", f.GetTrafficRule)
	r.GET("/zone/:name/portForward", f.GetPortForward)
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

func (f *Firewall) GetZoneNames(ctx *router.Context) error {
	names, err := firewall.GetZoneNames(permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(names)
}

func (f *Firewall) GetZoneByName(ctx *router.Context) error {
	zone, err := firewall.GetZone(ctx.Param("name"), permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(zone)
}

func (f *Firewall) GetZone(ctx *router.Context) error {
	zones, err := firewall.GetZones(permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(zones)
}

func (f *Firewall) UpdateZoneByName(ctx *router.Context) error {
	zone := &firewall.Zone{}
	if err := ctx.ShouldBindJSON(zone); err != nil {
		return err
	}

	err := firewall.UpdateZone(zone, permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) GetTrafficRule(ctx *router.Context) error {
	rules, err := firewall.GetTrafficRules(ctx.Param("name"), permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(rules)
}

func (f *Firewall) GetPortForward(ctx *router.Context) error {
	forwards, err := firewall.GetPortForwards(ctx.Param("name"), permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(forwards)
}
