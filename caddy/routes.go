package caddy

import (
	"errors"
	"fmt"
	"gitlab.com/xiayesuifeng/gopanel/caddy/config"
)

var (
	RouteNotFoundError       = errors.New("route not found")
	RouteGroupMustExistError = errors.New("route group must exist")
)

func GetRouteIdx(group string) (int, error) {
	routes := make([]config.RouteType, 0)

	resp, err := getClient().R().SetResult(routes).Get(getDefaultRouteApi())
	if err != nil {
		return -1, err
	}

	if resp.StatusCode() != 200 {
		return -1, errors.New("caddy admin api return message: " + string(resp.Body()))
	}

	for i, route := range routes {
		if route.Group == group {
			return i, nil
		}
	}

	return -1, RouteNotFoundError
}

func AddRoute(routeConfig *config.RouteType) error {
	if routeConfig.Group == "" {
		return RouteGroupMustExistError
	}

	resp, err := getClient().R().
		SetHeader("Content-Type", "application/json").
		SetBody(routeConfig).
		Post(getDefaultRouteApi())
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("caddy admin api return message: " + string(resp.Body()))
	} else {
		return nil
	}
}

func EditRoute(routeConfig *config.RouteType) error {
	if routeConfig.Group == "" {
		return RouteGroupMustExistError
	}

	idx, err := GetRouteIdx(routeConfig.Group)
	if err != nil {
		return err
	}

	resp, err := getClient().R().
		SetHeader("Content-Type", "application/json").
		SetBody(routeConfig).
		Post(fmt.Sprintf("%s/%d", getDefaultRouteApi(), idx))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("caddy admin api return message: " + string(resp.Body()))
	} else {
		return nil
	}
}

func DeleteRoute(group string) error {
	idx, err := GetRouteIdx(group)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%d", getDefaultRouteApi(), idx)

	resp, err := getClient().R().Delete(path)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("caddy admin api return message: " + string(resp.Body()))
	} else {
		return nil
	}
}

func getDefaultRouteApi() string {
	return fmt.Sprintf("%s/%s/routes", serversApi, DefaultHttpsServerName)
}
