package caddymodule

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"time"
)

type getOfficialPluginListResp struct {
	StatusCode int           `json:"status_code"`
	Result     []*PluginInfo `json:"result"`
}

type PluginInfo struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Published time.Time `json:"published"`
	Updated   time.Time `json:"updated"`
	Listed    bool      `json:"listed"`
	Available bool      `json:"available"`
	Downloads int       `json:"downloads"`
	Modules   []*Module `json:"modules"`
	Repo      string    `json:"repo"`
}

type Module struct {
	Name    string `json:"name"`
	Package string `json:"package"`
	Repo    string `json:"repo"`
}

func GetOfficialPluginList() ([]*PluginInfo, error) {
	resp := &getOfficialPluginListResp{}

	_, err := resty.New().
		R().
		SetResult(resp).
		Get("https://caddyserver.com/api/packages")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return resp.Result, nil
	} else {
		return nil, errors.New("return status_code not 200")
	}
}
