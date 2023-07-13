package network

import (
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
)

type Network struct {
}

func (n *Network) Name() string {
	return "network"
}

func (n *Network) Run(r router.Router) {
	r.GET("/devices", n.GetDevices)
}
