package caddyManager

import (
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyapp"
	"log"
	"strings"
)

type Config struct {
	Apps    map[string]interface{} `json:"apps"`
	Logging caddy.Logging          `json:"logging"`
}

func (m *Manager) convertToCaddyConfig() (config *Config) {
	m.filterInvalidDNSChallenges()
	m.filterInvalidWildcardDomains()

	config = newConfig()

	app := config.Apps["http"].(*caddyhttp.App)

	app.HTTPPort = m.caddyConf.General.HTTPPort
	app.HTTPSPort = m.caddyConf.General.HTTPSPort

	servers := app.Servers

	serverIdx := 1
	serverName := map[int]string{app.HTTPSPort: m.HTTPSServerName}

	servers[m.HTTPSServerName] = newServer(m.caddyConf.General.AllowH2C, fmt.Sprintf(":%d", app.HTTPSPort))

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
				servers[name] = newServer(m.caddyConf.General.AllowH2C, fmt.Sprintf(":%d", port))

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
			servers[name] = newServer(m.caddyConf.General.AllowH2C, fmt.Sprintf(":%d", config.ListenPort))

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

	// load caddy app config
	for name, app := range m.caddyApp {
		if !m.moduleList.HasNonStandardModule(name) {
			log.Println("[caddy manager] caddy module:", name, "not found,skip to loading non standard app")
			continue
		}

		if c := app.LoadConfig(caddyapp.Context{Change: m.appChange, ModuleList: m.moduleList}); c != nil {
			config.Apps[name] = c
		}

	}

	return
}

func (m *Manager) filterInvalidDNSChallenges() {
	for domain, config := range m.caddyConf.TLS.DNSChallenges {
		provider := make(map[string]string)

		err := json.Unmarshal(config.ProviderRaw, &provider)
		if err != nil {
			log.Println("[caddy manager] unmarshal dns challenges provider fail,error:", err)
			return
		}

		if !m.moduleList.HasNonStandardModule("dns.providers." + provider["name"]) {
			delete(m.caddyConf.TLS.DNSChallenges, domain)
		}
	}

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

func newServer(allowH2C bool, listen ...string) *caddyhttp.Server {
	srv := &caddyhttp.Server{
		Listen: listen,
	}

	if allowH2C {
		srv.Protocols = []string{"h1", "h2", "h2c", "h3"}
	}

	return srv
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
