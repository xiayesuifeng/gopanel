package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/service"
	"sort"
)

type Service struct {
}

func (s *Service) Get(ctx *gin.Context) {
	if services, err := service.GetServices(context.TODO()); err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	} else {
		sort.Sort(services)

		ctx.JSON(200, gin.H{
			"code": 200,
			"data": services,
		})
	}
}

func (s *Service) Post(ctx *gin.Context) {
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
	default:
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": "action must be one of start, stop, restart, enable",
		})
		return
	}

	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": 200,
			"data": data,
		})
	}
}
