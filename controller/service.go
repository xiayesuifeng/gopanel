package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/service"
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
		ctx.JSON(200, gin.H{
			"code": 200,
			"data": services,
		})
	}
}
