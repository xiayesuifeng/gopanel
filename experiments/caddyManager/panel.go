package caddyManager

import (
	"encoding/json"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"strings"
)

const panelRouteHandleCaddyJson = `{
      "handler": "subroute",
      "routes": [
        {netdataJson}
        {
          "handle": [{
			"handler": "reverse_proxy",
            "upstreams": [{"dial": "127.0.0.1:{port}"}]
          }]
        }
      ]
    }`

const netdataPathCaddyJson = `{
	"handler": "rewrite",
    "uri_substring": [{
	  "find": "/",
	  "replace": "{path}",
	  "limit": 1
	}]
  },`

const netdataCaddyJson = `{
      "handle": [{
        "handler": "static_response",
        "headers": {"Location": ["/netdata/"]},
        "status_code": 302
      }],
      "match": [{"path": ["/netdata"]}]
	},
	{
      "handle": [
	    {
	      "handler": "rewrite",
	      "strip_path_prefix": "/netdata"
	    },
	    {netdataPathJson}
        {
	      "handler": "reverse_proxy",
	      "upstreams": [{"dial": "{host}"}],
	      {netdataSSLJson}
          "headers": {"request": {"set": {"Host": ["{http.reverse_proxy.upstream.hostport}"]}}}
	    }
	  ],
	  "match": [{"path": ["/netdata/*"]}]
	},`

func loadPanelConfig(port string) *APPConfig {
	panelApp := &APPConfig{
		Domain:     []string{config.Conf.Panel.Domain},
		DisableSSL: config.Conf.Panel.DisableSSL,
	}

	if port := config.Conf.Panel.Port; port != 0 {
		panelApp.ListenPort = port
	}

	conf := panelRouteHandleCaddyJson
	conf = strings.ReplaceAll(conf, "{port}", port)

	if config.Conf.Netdata.Enable {
		netdataConf := netdataCaddyJson
		netdataConf = strings.ReplaceAll(netdataConf, "{host}", config.Conf.Netdata.Host)
		netdataPath := config.Conf.Netdata.Path
		if netdataPath != "" && netdataPath != "/" {
			if !strings.HasSuffix(netdataPath, "/") {
				netdataPath += "/"
			}
			pathConf := netdataPathCaddyJson
			pathConf = strings.ReplaceAll(pathConf, "{path}", netdataPath)
			netdataConf = strings.ReplaceAll(netdataConf, "{netdataPathJson}", pathConf)
		} else {
			netdataConf = strings.ReplaceAll(netdataConf, "{netdataPathJson}", "")
		}

		if config.Conf.Netdata.SSL {
			netdataConf = strings.ReplaceAll(netdataConf, "{netdataSSLJson}", "\"transport\": {\"protocol\": \"http\",\"tls\": {}},")
		} else {
			netdataConf = strings.ReplaceAll(netdataConf, "{netdataSSLJson}", "")
		}

		conf = strings.ReplaceAll(conf, "{netdataJson}", netdataConf)
	} else {
		conf = strings.ReplaceAll(conf, "{netdataJson}", "")
	}

	panelApp.Routes = append(panelApp.Routes, caddyhttp.Route{
		HandlersRaw: []json.RawMessage{
			json.RawMessage(conf),
		},
	})

	return panelApp
}
