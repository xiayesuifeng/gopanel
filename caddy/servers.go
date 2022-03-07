package caddy

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/xiayesuifeng/gopanel/caddy/config"
	"strings"
)

const serversApi = "/config/apps/http/servers"

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

func AddServer(name string, config config.ServerType) error {
	resp, err := getClient().R().
		SetHeader("Content-Type", "application/json").
		SetBody(&config).
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

func EditServer(name string, config config.ServerType) error {
	resp, err := getClient().R().
		SetHeader("Content-Type", "application/json").
		SetBody(&config).
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
		return strings.TrimSpace(string(resp.Body())) != "null"
	}
}

// Deprecated: use AddRoute instead.
func AddAdaptServerToRoute(name string, config config.ServerType) error {
	routes := config.Routes
	for i := 0; i < len(routes); i++ {
		routes[i].Group = name

		if err := AddRoute(routes[i]); err != nil {
			return err
		}
	}

	return nil
}

// Deprecated: use EditRoute instead.
func EditAdaptServerToRoute(name string, config config.ServerType) error {
	routes := config.Routes
	for i := 0; i < len(routes); i++ {
		routes[i].Group = name

		if err := EditRoute(routes[i]); err != nil {
			return err
		}
	}

	return nil
}

// Deprecated: use DeleteRoute instead.
func DeleteAdaptServerToRoute(name string) error {
	return DeleteRoute(name)
}
