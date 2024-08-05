package containify

import (
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/module/containify/engine"
	"strings"
)

func (c *Containify) GetImages(ctx *router.Context) error {
	list, err := c.service.ContainerEngine().Image().List(ctx, true, nil)
	if err != nil {
		return err
	}

	return ctx.JSON(list)
}

func (c *Containify) RemoveImage(ctx *router.Context) error {
	nameOrID := ctx.Param("nameOrID")

	if err := c.service.ContainerEngine().Image().Remove(ctx, nameOrID); err != nil {
		return err
	}

	return ctx.NoContent()
}

func (c *Containify) InspectImage(ctx *router.Context) error {
	name := strings.TrimPrefix(ctx.Param("name"), "/")

	image, err := engine.InspectImage(ctx, name)
	if err != nil {
		return err
	}

	return ctx.JSON(image)
}
