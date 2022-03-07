package caddy

import (
	"encoding/json"
	"fmt"
	"gitlab.com/xiayesuifeng/gopanel/caddy/config"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"strings"
)

const panelRouteHandleCaddyJson = `[
    {
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
    }
  ]`

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
	      "transport": {
	        "protocol": "http",
	    	"tls": {}
	      },
          "headers": {"request": {"set": {"Host": ["{http.reverse_proxy.upstream.hostport}"]}}}
	    }
	  ],
	  "match": [{"path": ["/netdata/*"]}]
	},`

func LoadPanelConfig(port string) (err error) {
	panelRoute := &config.RouteType{
		Group: "gopanel",
	}

	conf := panelRouteHandleCaddyJson
	conf = strings.ReplaceAll(conf, "{port}", port)
	if core.Conf.Panel.Domain != "" {
		panelRoute.Match = append(panelRoute.Match, map[string]json.RawMessage{
			"host": json.RawMessage("[\"" + core.Conf.Panel.Domain + "\"]"),
		})
	}

	if core.Conf.Netdata.Enable {
		netdataConf := netdataCaddyJson
		netdataConf = strings.ReplaceAll(netdataConf, "{host}", core.Conf.Netdata.Host)
		netdataPath := core.Conf.Netdata.Path
		if netdataPath != "" && netdataPath != "/" {
			if !strings.HasSuffix(netdataPath, "/") {
				netdataPath += "/"
			}
			pathConf := netdataPathCaddyJson
			pathConf = strings.ReplaceAll(pathConf, "{path}", netdataPath)
			netdataConf = strings.ReplaceAll(netdataConf, "{netdataPathJson}", pathConf)
		}
		conf = strings.ReplaceAll(conf, "{netdataJson}", netdataConf)
	} else {
		conf = strings.ReplaceAll(conf, "{netdataJson}", "")
	}

	panelRoute.Handle = json.RawMessage(conf)

	if CheckServerExist(DefaultHttpsServerName) {
		if _, err = GetRouteIdx(panelRoute.Group); err == RouteNotFoundError {
			err = AddRoute(panelRoute)
		} else {
			err = EditRoute(panelRoute)
		}
	} else {
		err = AddServer(DefaultHttpsServerName, config.ServerType{
			Listen: []string{
				fmt.Sprintf(":%d", DefaultHttpsPort),
			},
			Routes: []*config.RouteType{panelRoute},
		})
	}

	return AddTLSPolicy([]string{core.Conf.Panel.Domain})
}
