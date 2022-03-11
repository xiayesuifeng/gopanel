package caddyManager

import (
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func (m *Manager) convertToCaddyConfig() (servers map[string]*caddyhttp.Server) {
	servers = make(map[string]*caddyhttp.Server)

	serverIdx := 1
	serverName := map[int]string{m.HTTPSPort: m.HTTPSServerName}

	servers[m.HTTPSServerName] = newServer(fmt.Sprintf(":%d", m.HTTPSPort))

	for _, config := range m.app {
		if config.ListenPort == 0 {
			config.ListenPort = m.HTTPSPort
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

func newServer(listen ...string) *caddyhttp.Server {
	return &caddyhttp.Server{
		Listen: listen,
	}
}
