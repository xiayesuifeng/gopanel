package network

import (
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/module/network"
)

func (*Network) GetDevices(ctx *router.Context) error {
	devices, err := network.GetDevices()
	if err != nil {
		return err
	}

	return ctx.JSON(devices)
}
