package caddy

import (
	"encoding/json"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"strconv"
	"strings"
)

const panelCaddyJson = `{
    {listenPort}
    "routes": [
      {
        "handle": [
          {
            "handler": "subroute",
            "routes": [
              {netdataJson}{
                "handle": [
				  {
				    "handler": "reverse_proxy",
                    "upstreams": [{"dial": "127.0.0.1:{port}"}]
                  }
                ]
              }
            ]
          }
        ],
        "match": [{"host": [{domain}]}]
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
	conf := panelCaddyJson
	conf = strings.ReplaceAll(conf, "{port}", port)
	conf = strings.ReplaceAll(conf, "{domain}", "\""+core.Conf.Panel.Domain+"\"")
	if core.Conf.Panel.Port != 0 {
		conf = strings.ReplaceAll(conf, "{listenPort}", "\"listen\": [\":"+strconv.Itoa(core.Conf.Panel.Port)+"\"],")
	} else {
		conf = strings.ReplaceAll(conf, "{listenPort}", "")
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

	if CheckServerExist("gopanel") {
		err = EditServer("gopanel", json.RawMessage(conf))
	} else {
		err = AddServer("gopanel", json.RawMessage(conf))
	}

	return err
}
