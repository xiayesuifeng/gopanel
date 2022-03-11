package caddyManager

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/go-resty/resty/v2"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"net"
	"strconv"
)

var (
	manager *Manager
)

const (
	DefaultHttpPort        = 80
	DefaultHttpsPort       = 443
	DefaultHttpsServerName = "gopanel"
)

type APPConfig struct {
	ListenPort uint
	Domain     string
	routes     []*caddyhttp.Route
}

type Manager struct {
	httpClient      *resty.Client
	HTTPPort        int
	HTTPSPort       int
	HTTPSServerName string
	App             map[string]*APPConfig
}

func InitManager(adminAddress core.NetAddress) (err error) {
	manager = &Manager{
		httpClient:      newClient(adminAddress),
		HTTPPort:        DefaultHttpPort,
		HTTPSPort:       DefaultHttpsPort,
		HTTPSServerName: DefaultHttpsServerName,
	}

	appConfig, err := manager.getAppConfig()
	if err != nil {
		return
	}

	if appConfig.HTTPPort != 0 {
		manager.HTTPPort = appConfig.HTTPPort
	}

	if appConfig.HTTPSPort != 0 {
		manager.HTTPSPort = appConfig.HTTPSPort
	}

	for name, server := range appConfig.Servers {
		for _, listen := range server.Listen {
			if _, port, err := net.SplitHostPort(listen); err == nil && strconv.Itoa(DefaultHttpsPort) == port {
				manager.HTTPSServerName = name
			}
		}
	}

	return nil
}

func GetManager() *Manager {
	return manager
}
