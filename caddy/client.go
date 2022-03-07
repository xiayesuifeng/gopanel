package caddy

import (
	"context"
	"github.com/go-resty/resty/v2"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"net"
	"net/http"
	"sync"
)

var (
	once   = sync.Once{}
	client *resty.Client
)

func getClient() *resty.Client {
	once.Do(func() {
		address := core.Conf.Caddy.AdminAddress

		client = resty.New()

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
	})

	return client
}
