package containify

import "gitlab.com/xiayesuifeng/gopanel/api/server/router"

func (c *Containify) GetContainers(ctx *router.Context) error {
	list, err := c.service.ContainerEngine().Container().List(ctx)
	if err != nil {
		return err
	}

	return ctx.JSON(list)
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
