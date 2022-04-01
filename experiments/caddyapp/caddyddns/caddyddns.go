package caddyddns

import (
	"encoding/json"
	"github.com/mholt/caddy-dynamicdns"
	"gitlab.com/xiayesuifeng/gopanel/core/settingStorage"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyapp"
)

const module = "caddy"

var ctx caddyapp.Context

type CaddyDDNS struct {
	Enabled bool           `json:"enabled"`
	Config  dynamicdns.App `json:"config"`
}

func (c *CaddyDDNS) AppInfo() caddyapp.AppInfo {
	return caddyapp.AppInfo{
		Name: "dynamic_dns",
		New: func(c caddyapp.Context) caddyapp.CaddyApp {
			ctx = c

			return &CaddyDDNS{}
		},
	}
}

func (c *CaddyDDNS) LoadConfig(ctx caddyapp.Context) any {
	return nil
}

func SetCaddyDDNS(ddns *CaddyDDNS) error {
	bytes, err := json.Marshal(ddns)
	if err != nil {
		return err
	}

	storage := settingStorage.GetStorage()

	err = storage.Set(module, "ddns", bytes)
	if err != nil {
		return err
	}

	ctx.NotifyChange()

	return nil
}

func GetCaddyDDNS() (*CaddyDDNS, error) {
	storage := settingStorage.GetStorage()

	bytes := storage.Get(module, "ddns", []byte("{}"))

	ddns := &CaddyDDNS{}

	return ddns, json.Unmarshal(bytes, ddns)
}
