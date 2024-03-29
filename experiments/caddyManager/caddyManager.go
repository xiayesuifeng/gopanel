package caddyManager

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/go-resty/resty/v2"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyapp"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyapp/caddyddns"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyutil/caddyconfig"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyutil/caddymodule"
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
	caddyConf       *caddyconfig.Configuration
	HTTPSServerName string
	app             map[string]*APPConfig
	appChange       chan bool
	onAppChange     func()

	moduleList *caddymodule.ModuleList

	caddyApp map[string]caddyapp.CaddyApp
	appMutex sync.RWMutex
}

func InitManager(adminAddress config.NetAddress, panelPort string) (err error) {
	manager = &Manager{
		httpClient:      newClient(adminAddress),
		HTTPSServerName: DefaultHttpsServerName,
		app:             map[string]*APPConfig{},
		appChange:       make(chan bool),
		caddyApp:        map[string]caddyapp.CaddyApp{},
	}

	manager.moduleList, err = caddymodule.GetModuleList()
	if err != nil {
		log.Println("[caddy manager] get caddy module list fail,error:", err)
		return
	}

	appConfig, err := manager.getAppConfig()
	if err != nil {
		return
	}

	httpPort := config.Conf.Caddy.DefaultHTTPPort
	httpsPort := config.Conf.Caddy.DefaultHTTPSPort

	if appConfig.HTTPPort != 0 {
		httpPort = appConfig.HTTPPort
	}

	if appConfig.HTTPSPort != 0 {
		httpsPort = appConfig.HTTPSPort
	}

	if err := caddyconfig.InitDefaultPortConf(httpPort, httpsPort); err != nil {
		return err
	}

	manager.caddyConf = caddyconfig.GetConfiguration()

	for name, server := range appConfig.Servers {
		for _, listen := range server.Listen {
			if _, port, err := net.SplitHostPort(listen); err == nil && strconv.Itoa(manager.caddyConf.General.HTTPSPort) == port {
				manager.HTTPSServerName = name
			}
		}
	}

	manager.RegisterCaddyApp(&caddyddns.CaddyDDNS{})

	manager.onAppChange = manager.onAppChangeFunc
	go manager.onAppChange()

	manager.AddOrUpdateApp("gopanel", loadPanelConfig(panelPort))

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

	m.caddyConf = caddyconfig.GetConfiguration()

	m.appChange <- true
}

func (m *Manager) RegisterCaddyApp(app caddyapp.CaddyApp) {
	m.appMutex.Lock()
	defer m.appMutex.Unlock()

	a := app.AppInfo()
	if a.Name == "" {
		panic("CaddyAPP Name missing")
	}

	m.caddyApp[a.Name] = a.New(caddyapp.Context{Change: m.appChange})
}
