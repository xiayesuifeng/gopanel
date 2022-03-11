package caddyManager

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/go-resty/resty/v2"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"log"
	"net"
	"strconv"
	"sync"
)

var (
	manager *Manager
	mutex   sync.Mutex
)

const (
	DefaultHttpPort        = 80
	DefaultHttpsPort       = 443
	DefaultHttpsServerName = "gopanel"
)

type APPConfig struct {
	ListenPort int               `json:"listenPort,omitempty"`
	Domain     []string          `json:"domain"`
	Routes     []caddyhttp.Route `json:"routes"`
}

type Manager struct {
	httpClient      *resty.Client
	HTTPPort        int
	HTTPSPort       int
	HTTPSServerName string
	app             map[string]*APPConfig
	appChange       chan bool
	onAppChange     func()
}

func InitManager(adminAddress core.NetAddress) (err error) {
	manager = &Manager{
		httpClient:      newClient(adminAddress),
		HTTPPort:        DefaultHttpPort,
		HTTPSPort:       DefaultHttpsPort,
		HTTPSServerName: DefaultHttpsServerName,
		app:             map[string]*APPConfig{},
		appChange:       make(chan bool),
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

	manager.onAppChange = manager.onAppChangeFunc
	go manager.onAppChange()

	return nil
}

func GetManager() *Manager {
	return manager
}

func (m *Manager) IsAppExist(name string) bool {
	_, exist := m.app[name]
	return exist
}

func (m *Manager) AddOrUpdateApp(name string, config *APPConfig) error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Println("[caddy manager] add or update app:", name)

	m.app[name] = config

	m.appChange <- true

	return nil
}

func (m *Manager) DeleteApp(name string) error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Println("[caddy manager] delete app:", name)

	delete(m.app, name)

	m.appChange <- true

	return nil
}

func (m *Manager) onAppChangeFunc() {
	log.Println("[caddy manager] start listening for app change")

	for <-m.appChange {
		log.Println("[caddy manager] app changes detected, call caddy admin api to update routes json")

		servers := m.convertToCaddyConfig()

		if err := m.postCaddyObject("/config/apps/http/servers", servers); err != nil {
			log.Println("[caddy manager] call caddy admin api fail: ", err)
		}
	}

	log.Println("[caddy manager] stop listening for app change")
}
