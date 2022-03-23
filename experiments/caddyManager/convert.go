package caddyManager

import (
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	"strings"
)

type Config struct {
	Apps    map[string]interface{} `json:"apps"`
	Logging caddy.Logging          `json:"logging"`
}

func (m *Manager) convertToCaddyConfig() (config *Config) {
	m.filterInvalidWildcardDomains()

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

	wildcardDomainsApp, normalApp := filterWildcardDomainsApp(m.caddyConf.TLS.WildcardDomains, m.app)

	var tlsDomains []string

	for wildcardDomain, configs := range wildcardDomainsApp {
		routesMap := map[int][]caddyhttp.Route{}

		for _, config := range configs {
			if config.ListenPort == 0 {
				config.ListenPort = app.HTTPSPort
			}

			route := newRoute(config.Domain, config.Routes)

			routesMap[config.ListenPort] = append(routesMap[config.ListenPort], route)
		}

		for port, routes := range routesMap {
			name, exist := serverName[port]

			if !exist {
				name = fmt.Sprintf("gopanel%d", serverIdx)
				serverName[port] = name
				servers[name] = newServer(fmt.Sprintf(":%d", port))

				serverIdx++
			}

			route := newRoute([]string{wildcardDomain}, routes)

			servers[name].Routes = append(servers[name].Routes, route)
		}

		tlsDomains = append(tlsDomains, wildcardDomain)
	}

	for _, config := range normalApp {
		if config.ListenPort == 0 {
			if config.DisableSSL {
				config.ListenPort = app.HTTPPort
			} else {
				config.ListenPort = app.HTTPSPort
			}
		}

		name, exist := serverName[config.ListenPort]

		if !exist {
			name = fmt.Sprintf("gopanel%d", serverIdx)
			serverName[config.ListenPort] = name
			servers[name] = newServer(fmt.Sprintf(":%d", config.ListenPort))

			serverIdx++
		}

		route := newRoute(config.Domain, config.Routes)

		servers[name].Routes = append(servers[name].Routes, route)

		if config.ListenPort != app.HTTPPort && config.DisableSSL {
			servers[name].AutoHTTPS = &caddyhttp.AutoHTTPSConfig{
				Disabled: true,
			}
		}

		tlsDomains = append(tlsDomains, config.Domain...)
	}

	config.Apps["tls"] = loadTLSConfig(tlsDomains, m.caddyConf.TLS.DNSChallenges)

	return
}

func (m *Manager) filterInvalidWildcardDomains() {
	var wildcardDomains []string

	tls := m.caddyConf.TLS
	for _, wildcardDomain := range tls.WildcardDomains {
		for domain := range tls.DNSChallenges {
			if strings.HasSuffix(wildcardDomain, domain) {
				wildcardDomains = append(wildcardDomains, wildcardDomain)
				break
			}
		}
	}

	m.caddyConf.TLS.WildcardDomains = wildcardDomains
}

func newRoute(domains []string, routes []caddyhttp.Route) caddyhttp.Route {
	match := caddy.ModuleMap{}

	if len(domains) != 0 {
		match["host"] = caddyconfig.JSON(domains, nil)
	}

	route := caddyhttp.Route{
		MatcherSetsRaw: []caddy.ModuleMap{match},
		HandlersRaw:    []json.RawMessage{caddyconfig.JSON(NewSubrouteHandle(routes), nil)},
		Terminal:       true,
	}

	return route
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

func filterWildcardDomainsApp(wildcardDomains []string, apps map[string]*APPConfig) (wildcardDomainAPPConfig map[string][]*APPConfig, normalAPPConfig []*APPConfig) {
	wildcardDomainAPPConfig = map[string][]*APPConfig{}

	for _, c := range apps {
		for _, domain := range c.Domain {
			wildcardDomain := findWildcardDomain(domain, wildcardDomains)
			if wildcardDomain != "" {
				wildcardDomainAPPConfig[wildcardDomain] = append(wildcardDomainAPPConfig[wildcardDomain], c)
			} else {
				normalAPPConfig = append(normalAPPConfig, c)
			}
		}
	}

	return
}

func findWildcardDomain(domain string, wildcardDomains []string) string {
	strs := strings.Split(domain, ".")

	if len(strs) < 2 {
		return ""
	}

	strs[0] = "*"
	targetDomain := strings.Join(strs, ".")

	for _, wildcardDomain := range wildcardDomains {
		if wildcardDomain == targetDomain {
			return wildcardDomain
		}
	}

	return ""
}
