package firewall

import (
	"errors"
	"gitlab.com/xiayesuifeng/gopanel/api/server/middleware"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/module/firewall"
)

type Firewall struct {
}

func (f *Firewall) Name() string {
	return "firewall"
}

func (f *Firewall) Run(r router.Router) {
	r.Use(middleware.BinaryExistMiddleware("firewalld"))
	r.Use(permanentMiddleware)

	r.GET("", f.GetConfig)
	r.POST("", f.UpdateConfig)

	r.POST("/reload", f.Reload)
	r.POST("/reset", f.Reset)

	r.GET("/zone", f.GetZone)
	r.POST("/zone", f.AddZone)

	r.GET("/zone/names", f.GetZoneNames)
	r.GET("/zone/:name", f.GetZoneByName)
	r.PUT("/zone/:name", f.UpdateZoneByName)
	r.DELETE("/zone/:name", f.RemoveZoneByName)

	r.GET("/zone/:name/trafficRule", f.GetTrafficRule)
	r.POST("/zone/:name/trafficRule", f.AddTrafficRule)
	r.DELETE("/zone/:name/trafficRule", f.RemoveTrafficRule)

	r.GET("/zone/:name/portForward", f.GetPortForward)
	r.POST("/zone/:name/portForward", f.AddPortForward)
	r.DELETE("/zone/:name/portForward", f.RemovePortForward)

	r.GET("/service/names", f.GetServiceNames)

	r.GET("/icmpType/names", f.GetICMPTypeNames)

	r.GET("/policy", f.GetPolicy)
	r.POST("/policy", f.AddPolicy)
	r.PUT("/policy/:name", f.UpdatePolicyByName)
}

type Config struct {
	DefaultZone string `json:"defaultZone" binding:"required"`
}

func (f *Firewall) GetConfig(ctx *router.Context) error {
	zone, err := firewall.GetDefaultZone()
	if err != nil {
		return err
	}

	return ctx.JSON(&Config{DefaultZone: zone})
}

func (f *Firewall) UpdateConfig(ctx *router.Context) error {
	config := &Config{}
	if err := ctx.ShouldBind(config); err != nil {
		return err
	}

	if err := firewall.SetDefaultZone(config.DefaultZone); err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) Reload(ctx *router.Context) error {
	if err := firewall.Reload(); err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) Reset(ctx *router.Context) error {
	if err := firewall.Reset(); err != nil {
		return err
	}

	return ctx.NoContent()
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

func (f *Firewall) AddZone(ctx *router.Context) error {
	zone := &firewall.Zone{}
	if err := ctx.ShouldBindJSON(zone); err != nil {
		return err
	}

	err := firewall.AddZone(zone)
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) UpdateZoneByName(ctx *router.Context) error {
	zone := &firewall.Zone{}
	if err := ctx.ShouldBindJSON(zone); err != nil {
		return err
	}

	err := firewall.UpdateZone(ctx.Param("name"), zone, permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) RemoveZoneByName(ctx *router.Context) error {
	if err := firewall.RemoveZone(ctx.Param("name")); err != nil {
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

type TrafficRuleRequest struct {
	firewall.TrafficRule

	Type firewall.RuleType `json:"type" binding:"required"`
}

func (f *Firewall) AddTrafficRule(ctx *router.Context) error {
	request := &TrafficRuleRequest{}

	if err := ctx.ShouldBind(request); err != nil {
		return ctx.Error(400, err.Error())
	}

	request.TrafficRule.Type = request.Type

	err := firewall.AddTrafficRule(ctx.Param("name"), &request.TrafficRule, permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) RemoveTrafficRule(ctx *router.Context) error {
	request := &TrafficRuleRequest{}

	if err := ctx.ShouldBind(request); err != nil {
		return ctx.Error(400, err.Error())
	}

	request.TrafficRule.Type = request.Type

	err := firewall.RemoveTrafficRule(ctx.Param("name"), &request.TrafficRule, permanent(ctx))
	if err != nil {
		if errors.Is(err, firewall.NotFoundErr) {
			return ctx.Error(404, "traffic rule not found")
		} else {
			return err
		}
	}

	return ctx.NoContent()
}

type PortForwardRequest struct {
	// Port port number or range
	Port     string                   `json:"port" binding:"required"`
	Protocol firewall.ForwardProtocol `json:"protocol" binding:"required"`
	// ToPort port number or range
	ToPort    string `json:"toPort" binding:"required"`
	ToAddress string `json:"toAddress"`
}

func (f *Firewall) AddPortForward(ctx *router.Context) error {
	request := &PortForwardRequest{}

	if err := ctx.ShouldBind(request); err != nil {
		return ctx.Error(400, err.Error())
	}

	err := firewall.AddPortForward(ctx.Param("name"), &firewall.PortForward{
		Port:      request.Port,
		Protocol:  request.Protocol,
		ToPort:    request.ToPort,
		ToAddress: request.ToAddress,
	}, permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) RemovePortForward(ctx *router.Context) error {
	request := &PortForwardRequest{}

	if err := ctx.ShouldBind(request); err != nil {
		return ctx.Error(400, err.Error())
	}

	err := firewall.RemovePortForward(ctx.Param("name"), &firewall.PortForward{
		Port:      request.Port,
		Protocol:  request.Protocol,
		ToPort:    request.ToPort,
		ToAddress: request.ToAddress,
	}, permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) GetServiceNames(ctx *router.Context) error {
	names, err := firewall.GetServiceNames(permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(names)
}

func (f *Firewall) GetPolicy(ctx *router.Context) error {
	policies, err := firewall.GetPolicies(permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(policies)
}

func (f *Firewall) AddPolicy(ctx *router.Context) error {
	policy := firewall.Policy{}
	if err := ctx.ShouldBindJSON(&policy); err != nil {
		return err
	}

	err := firewall.AddPolicy(policy)
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) UpdatePolicyByName(ctx *router.Context) error {
	policy := firewall.Policy{}
	if err := ctx.ShouldBindJSON(&policy); err != nil {
		return err
	}

	err := firewall.UpdatePolicy(ctx.Param("name"), policy, permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (f *Firewall) GetICMPTypeNames(ctx *router.Context) error {
	names, err := firewall.GetICMPTypeNames(permanent(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(names)
}
