package config

import "encoding/json"

type HttpType struct {
	HttpPort    int                   `json:"http_port,omitempty"`
	HttpsPort   int                   `json:"https_port,omitempty"`
	GracePeriod int                   `json:"grace_period,omitempty"`
	Servers     map[string]ServerType `json:"servers,omitempty"`
}

type ServerType struct {
	Listen                []string            `json:"listen,omitempty"`
	ReadTimeout           int                 `json:"read_timeout,omitempty"`
	ReadHeaderTimeout     int                 `json:"read_header_timeout,omitempty"`
	WriteTimeout          int                 `json:"write_timeout,omitempty"`
	IdleTimeout           int                 `json:"idle_timeout,omitempty"`
	MaxHeaderBytes        int                 `json:"max_header_bytes,omitempty"`
	Routes                []*RouteType        `json:"routes,omitempty"`
	Errors                json.RawMessage     `json:"errors,omitempty"`
	TLSConnectionPolicies json.RawMessage     `json:"tls_connection_policies,omitempty"`
	AutomaticHttps        *AutomaticHttpsType `json:"automatic_https,omitempty"`
	StrictSniHost         bool                `json:"strict_sni_host,omitempty"`
	Logs                  json.RawMessage     `json:"logs,omitempty"`
	ExperimentalHttp3     bool                `json:"experimental_http3,omitempty"`
	AllowH2C              bool                `json:"allow_h2c,omitempty"`
}

type RouteType struct {
	Group    string                       `json:"group,omitempty"`
	Match    []map[string]json.RawMessage `json:"match,omitempty"`
	Handle   json.RawMessage              `json:"handle,omitempty"`
	Terminal bool                         `json:"terminal,omitempty"`
}

type AutomaticHttpsType struct {
	Disable                  bool     `json:"disable,omitempty"`
	DisableRedirects         bool     `json:"disable_redirects,omitempty"`
	Skip                     []string `json:"skip,omitempty"`
	SkipCertificates         []string `json:"skip_certificates,omitempty"`
	IgnoreLoadedCertificates bool     `json:"ignore_loaded_certificates,omitempty"`
}
