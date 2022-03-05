package caddy

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"net"
	"net/http"
	"sync"
)

const serversApi = "/config/apps/http/servers"

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

func GetServers() (json.RawMessage, error) {
	resp, err := getClient().R().Get(serversApi)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New("caddy admin api return message: " + string(resp.Body()))
	} else {
		return resp.Body(), nil
	}
}

func GetServer(name string) (json.RawMessage, error) {
	resp, err := getClient().R().Get(serversApi + "/" + name)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New("caddy admin api return message: " + string(resp.Body()))
	} else {
		return resp.Body(), nil
	}
}

func AddServer(name string, config json.RawMessage) error {
	resp, err := getClient().R().
		SetHeader("Content-Type", "application/json").
		SetBody(config).
		Put(serversApi + "/" + name)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("caddy admin api return message: " + string(resp.Body()))
	} else {
		return nil
	}
}

func EditServer(name string, config json.RawMessage) error {
	resp, err := getClient().R().
		SetHeader("Content-Type", "application/json").
		SetBody(config).
		Post(serversApi + "/" + name)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("caddy admin api return message: " + string(resp.Body()))
	} else {
		return nil
	}
}

func DeleteServer(name string) error {
	resp, err := getClient().R().Delete(serversApi + "/" + name)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("caddy admin api return message: " + string(resp.Body()))
	} else {
		return nil
	}
}

func CheckServerExist(name string) bool {
	resp, err := getClient().R().Get(serversApi + "/" + name)
	if err != nil {
		return false
	}

	if resp.StatusCode() != 200 {
		return false
	} else {
		return string(resp.Body()) != "null"
	}
}
