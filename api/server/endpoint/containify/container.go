package containify

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine/entity"
)

func (c *Containify) GetContainers(ctx *router.Context) error {
	list, err := c.service.ContainerEngine().Container().List(ctx)
	if err != nil {
		return err
	}

	return ctx.JSON(list)
}

type CreateContainer struct {
	Name       string               `json:"name" binding:"required"`
	Image      string               `json:"image" binding:"required"`
	Entrypoint string               `json:"entrypoint"`
	Command    []string             `json:"command"`
	Env        map[string]string    `json:"env"`
	Ports      []entity.PortMapping `json:"ports"`
	Labels     map[string]string    `json:"labels"`
	Mounts     []entity.Mount       `json:"mounts"`
}

func (c *Containify) CreateContainers(ctx *router.Context) error {
	data := &CreateContainer{}

	if err := ctx.ShouldBind(data); err != nil {
		return ctx.Error(400, err.Error())
	}

	exist, err := c.service.ContainerEngine().Image().Exists(ctx, data.Image)
	if err != nil {
		return err
	}

	if !exist {
		_, err := c.service.ContainerEngine().Image().Pull(ctx, data.Image, nil)
		if err != nil {
			return err
		}
	}

	id, err := c.service.ContainerEngine().Container().Create(ctx, &entity.Container{
		ContainerBasic: entity.ContainerBasic{
			Name:    data.Name,
			Image:   data.Image,
			Command: data.Command,
			Ports:   data.Ports,
			Env:     data.Env,
			Labels:  data.Labels,
		},
		Entrypoint: data.Entrypoint,
		Mounts:     data.Mounts,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(gin.H{
		"containerID": id,
	})
}

func (c *Containify) StartContainer(ctx *router.Context) error {
	nameOrID := ctx.Param("nameOrID")

	err := c.service.ContainerEngine().Container().Start(ctx, nameOrID)
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (c *Containify) RestartContainer(ctx *router.Context) error {
	nameOrID := ctx.Param("nameOrID")

	err := c.service.ContainerEngine().Container().Restart(ctx, nameOrID)
	if err != nil {
		return err
	}

	return ctx.NoContent()
}

func (c *Containify) StopContainer(ctx *router.Context) error {
	nameOrID := ctx.Param("nameOrID")

	err := c.service.ContainerEngine().Container().Stop(ctx, nameOrID)
	if err != nil {
		return err
	}

	return ctx.NoContent()
}
