package caddyddns

import (
	"encoding/json"
	"errors"
	"github.com/mholt/caddy-dynamicdns"
	"gitlab.com/xiayesuifeng/gopanel/core/settingStorage"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyapp"
	"log"
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
	ddns, err := GetCaddyDDNS()
	if err != nil {
		log.Println("[caddy ddns] get ddns config fail,err:", err)
		return nil
	}

	provider := make(map[string]string)

	err = json.Unmarshal(ddns.Config.DNSProviderRaw, &provider)
	if err != nil {
		log.Println("[caddy ddns] unmarshal dns provider fail,error:", err)
		return nil
	}

	moduleName := "dns.providers." + provider["name"]
	if !ctx.ModuleList.HasNonStandardModule(moduleName) {
		log.Println("[caddy ddns] caddy module:", moduleName, "not found,skip")
		return nil
	}

	if ddns.Enabled {
		return ddns.Config
	} else {
		return nil
	}
}

func SetCaddyDDNS(ddns *CaddyDDNS) error {
	if len(ddns.Config.DNSProviderRaw) == 0 {
		return errors.New("a DNS provider is required")
	}

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
