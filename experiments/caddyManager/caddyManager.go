package caddyManager

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/go-resty/resty/v2"
	"gitlab.com/xiayesuifeng/gopanel/configuration/caddy"
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

	DisableSSL bool `json:"disableSSL,omitempty"`
}

type Manager struct {
	httpClient      *resty.Client
	caddyConf       *caddy.Configuration
	HTTPSServerName string
	app             map[string]*APPConfig
	appChange       chan bool
	onAppChange     func()
}

func InitManager(adminAddress core.NetAddress) (err error) {
	manager = &Manager{
		httpClient:      newClient(adminAddress),
		HTTPSServerName: DefaultHttpsServerName,
		app:             map[string]*APPConfig{},
		appChange:       make(chan bool),
	}

	appConfig, err := manager.getAppConfig()
	if err != nil {
		return
	}

	httpPort := DefaultHttpPort
	httpsPort := DefaultHttpsPort

	if appConfig.HTTPPort != 0 {
		httpPort = appConfig.HTTPPort
	}

	if appConfig.HTTPSPort != 0 {
		httpsPort = appConfig.HTTPSPort
	}

	if err := caddy.InitDefaultPortConf(httpPort, httpsPort); err != nil {
		return err
	}

	manager.caddyConf = caddy.GetConfiguration()

	for name, server := range appConfig.Servers {
		for _, listen := range server.Listen {
			if _, port, err := net.SplitHostPort(listen); err == nil && strconv.Itoa(manager.caddyConf.General.HTTPSPort) == port {
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

		config := m.convertToCaddyConfig()

		if err := m.postCaddyObject("/config/apps", config.Apps); err != nil {
			log.Println("[caddy manager] call caddy admin api fail: ", err)
		}
	}

	log.Println("[caddy manager] stop listening for app change")
}

func (m *Manager) NotifyCaddyConfigChange() {
	oldHttpsPort := m.caddyConf.General.HTTPSPort

	for name := range m.app {
		if m.app[name].ListenPort == oldHttpsPort {
			m.app[name].ListenPort = 0
		}
	}

	m.caddyConf = caddy.GetConfiguration()

	m.appChange <- true
}
