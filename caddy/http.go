package caddy

import (
	"errors"
	"gitlab.com/xiayesuifeng/gopanel/caddy/config"
	"net"
	"strconv"
)

const httpApi = "/config/apps/http"

var (
	DefaultHttpPort        = 80
	DefaultHttpsPort       = 443
	DefaultHttpsServerName = "gopanel"
)

func GetHttpConfig() (httpConfig *config.HttpType, err error) {
	httpConfig = &config.HttpType{}

	resp, err := getClient().R().SetResult(httpConfig).Get(httpApi)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New("caddy admin api return message: " + string(resp.Body()))
	}

	return
}

func ParseDefaultHTTPConfig() error {
	httpConfig, err := GetHttpConfig()
	if err != nil {
		return err
	}

	if httpConfig.HttpPort != 0 {
		DefaultHttpPort = httpConfig.HttpPort
	}

	if httpConfig.HttpsPort != 0 {
		DefaultHttpsPort = httpConfig.HttpsPort
	}

	for name, server := range httpConfig.Servers {
		for _, listen := range server.Listen {
			if _, port, err := net.SplitHostPort(listen); err == nil && strconv.Itoa(DefaultHttpsPort) == port {
				DefaultHttpsServerName = name
			}
		}
	}

	return nil
}
