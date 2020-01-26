package caddy

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"gitlab.com/xiayesuifeng/gopanel/core"
)

const serversApi = "/config/apps/http/servers"

func GetServers() (json.RawMessage, error) {
	resp, err := resty.New().R().Get(core.Conf.Caddy.AdminAddress + serversApi)
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
	resp, err := resty.New().R().Get(core.Conf.Caddy.AdminAddress + serversApi + "/" + name)
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
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(config).
		Put(core.Conf.Caddy.AdminAddress + serversApi + "/" + name)
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
	resp, err := resty.New().R().Delete(core.Conf.Caddy.AdminAddress + serversApi + "/" + name)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("caddy admin api return message: " + string(resp.Body()))
	} else {
		return nil
	}
}
