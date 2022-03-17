package caddyManager

import (
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
)

type Config struct {
	Apps    map[string]interface{} `json:"apps"`
	Logging caddy.Logging          `json:"logging"`
}

func (m *Manager) convertToCaddyConfig() (config *Config) {
	config = newConfig()

	app := config.Apps["http"].(*caddyhttp.App)

	app.HTTPPort = m.caddyConf.General.HTTPPort
	app.HTTPSPort = m.caddyConf.General.HTTPSPort

	servers := app.Servers

	serverIdx := 1
	serverName := map[int]string{app.HTTPSPort: m.HTTPSServerName}

	servers[m.HTTPSServerName] = newServer(fmt.Sprintf(":%d", app.HTTPSPort))
	servers[m.HTTPSServerName].ExperimentalHTTP3 = m.caddyConf.General.ExperimentalHttp3
	servers[m.HTTPSServerName].AllowH2C = m.caddyConf.General.AllowH2C

	for _, config := range m.app {
		if config.ListenPort == 0 {
			config.ListenPort = app.HTTPSPort
		}

		name, exist := serverName[config.ListenPort]

		if !exist {
			name = fmt.Sprintf("gopanel%d", serverIdx)
			serverName[config.ListenPort] = name
			servers[name] = newServer(fmt.Sprintf(":%d", config.ListenPort))

			serverIdx++
		}

		match := caddy.ModuleMap{}

		if len(config.Domain) != 0 {
			match["host"] = caddyconfig.JSON(config.Domain, nil)
		}

		route := caddyhttp.Route{
			MatcherSetsRaw: []caddy.ModuleMap{match},
			HandlersRaw:    []json.RawMessage{caddyconfig.JSON(NewSubrouteHandle(config.Routes), nil)},
			Terminal:       true,
		}

		servers[name].Routes = append(servers[m.HTTPSServerName].Routes, route)
	}

	return
}

func newConfig() *Config {
	return &Config{
		Apps: map[string]interface{}{
			"http": &caddyhttp.App{
				Servers: map[string]*caddyhttp.Server{},
			},
			"tls": &caddytls.TLS{},
		},
	}
}

func newServer(listen ...string) *caddyhttp.Server {
	return &caddyhttp.Server{
		Listen: listen,
	}
}
