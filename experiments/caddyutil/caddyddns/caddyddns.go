package caddyddns

import (
	"encoding/json"
	"github.com/mholt/caddy-dynamicdns"
	"gitlab.com/xiayesuifeng/gopanel/core/settingStorage"
)

const module = "caddy"

type CaddyDDNS struct {
	Enabled bool           `json:"enabled"`
	Config  dynamicdns.App `json:"config"`
}

func SetCaddyDDNS(ddns *CaddyDDNS) error {
	bytes, err := json.Marshal(ddns)
	if err != nil {
		return err
	}

	storage := settingStorage.GetStorage()

	return storage.Set(module, "ddns", bytes)
}

func GetCaddyDDNS() (*CaddyDDNS, error) {
	storage := settingStorage.GetStorage()

	bytes := storage.Get(module, "ddns", []byte("{}"))

	ddns := &CaddyDDNS{}

	return ddns, json.Unmarshal(bytes, ddns)
}
