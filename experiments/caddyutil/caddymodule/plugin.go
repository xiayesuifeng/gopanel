package caddymodule

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"os/exec"
	"strings"
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

func InstallPlugin(packages ...string) error {
	path, err := exec.LookPath("caddy")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "add-package", strings.Join(packages, " "))

	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		str := string(outBytes)
		for _, line := range strings.Split(string(outBytes), "\n") {
			if strings.HasPrefix(line, "Error:") {
				return errors.New(strings.Replace(line, "Error: ", "", 1))
			}
		}

		return errors.New(str)
	}

	return nil
}
