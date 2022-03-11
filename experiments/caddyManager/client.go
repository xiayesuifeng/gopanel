package caddyManager

import (
	"context"
	"github.com/go-resty/resty/v2"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"net"
	"net/http"
	"time"
)

func newClient(address core.NetAddress) *resty.Client {
	client := resty.New()

	if address.IsUnixNetwork() {
		transport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial(address.Network, address.Address)
			},
		}

		client.SetTransport(transport).SetScheme("http")
	} else {
		client.SetHostURL(address.Address)
	}

	client.SetTimeout(time.Second * 10)
	client.SetRetryCount(3)

	return client
}
