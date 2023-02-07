package service

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/service"
	"sort"
)

type Service struct {
}

func (s *Service) Name() string {
	return "service"
}

func (s *Service) Run(r router.Router) {
	r.GET("", s.Get)
	r.POST("/:name/:action", s.Post)
}

func (s *Service) Get(ctx *router.Context) error {
	if services, err := service.GetServices(ctx); err != nil {
		return err
	} else {
		sort.Sort(services)

		return ctx.JSON(gin.H{
			"data": services,
		})
	}
}

func (s *Service) Post(ctx *router.Context) error {
	name := ctx.Param("name")
	action := ctx.Param("action")

	stopTriggeredBy := ctx.Query("stopTriggeredBy")

	var (
		data interface{}
		err  error
	)

	switch action {
	case "start":
		data, err = service.StartService(ctx, name, service.FailMode)
	case "stop":
		data, err = service.StopService(ctx, name, service.FailMode, stopTriggeredBy == "true")
	case "restart":
		data, err = service.RestartService(ctx, name, service.FailMode)
	case "enable":
		_, data, err = service.EnableService(ctx, name)
	case "disable":
		data, err = service.DisableService(ctx, name)
	default:
		return ctx.Error(400, "action must be one of start, stop, restart, enable, disable")
	}

	if err != nil {
		return err
	} else {
		return ctx.JSON(gin.H{
			"data": data,
		})
	}
}
