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
              {
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

func LoadPanelConfig(port string) (err error) {
	conf := panelCaddyJson
	conf = strings.ReplaceAll(conf, "{port}", port)
	conf = strings.ReplaceAll(conf, "{domain}", "\""+core.Conf.Panel.Domain+"\"")
	if core.Conf.Panel.Port != 0 {
		conf = strings.ReplaceAll(conf, "{listenPort}", "\"listen\": [\":"+strconv.Itoa(core.Conf.Panel.Port)+"\"],")
	} else {
		conf = strings.ReplaceAll(conf, "{listenPort}", "")
	}

	if CheckServerExist("gopanel") {
		err = EditServer("gopanel", json.RawMessage(conf))
	} else {
		err = AddServer("gopanel", json.RawMessage(conf))
	}

	return err
}
