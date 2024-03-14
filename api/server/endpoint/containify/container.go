package containify

import "gitlab.com/xiayesuifeng/gopanel/api/server/router"

func (c *Containify) GetContainers(ctx *router.Context) error {
	list, err := c.service.ContainerEngine().Container().List(ctx)
	if err != nil {
		return err
	}

	return ctx.JSON(list)
}
