package app

import (
	"encoding/json"
	"errors"
	"gitlab.com/xiayesuifeng/gopanel/backend"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyutil/caddyvalidate"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Type int

const (
	GoType Type = iota + 1
	JavaType
	PHPType
	OtherType
	ReverseProxyType
	StaticWebType
)

type BackendType string

const (
	NoneBackend BackendType = "none"
	ExecBackend BackendType = "exec"
)

type App struct {
	Name             string                 `json:"name"`
	CaddyConfig      caddyManager.APPConfig `json:"caddyConfig" binding:"required"`
	Type             Type                   `json:"type" binding:"required"`
	BackendType      BackendType            `json:"backendType"`
	BackendConfigRaw json.RawMessage        `json:"backendConfig"`
	Version          string                 `json:"version"`

	// Deprecated: use ExecBackendConfig instead. To be removed in 1.0.0 release.
	AutoReboot bool `json:"autoReboot,omitempty"`
	// Deprecated: use ExecBackendConfig instead. To be removed in 1.0.0 release.
	Path string `json:"path,omitempty"`
	// Deprecated: use ExecBackendConfig instead. To be removed in 1.0.0 release.
	BootArgument string `json:"bootArgument,omitempty"`
}

type ExecBackendConfig struct {
	WorkingDirectory string `json:"workingDirectory"`
	AutoReboot       bool   `json:"autoReboot"`
	Path             string `json:"path"`
	Argument         string `json:"argument"`
}

func AddApp(app App, validate bool) error {
	if CheckAppExist(app.Name) {
		return errors.New("app is exist")
	}

	if app.Version == "" {
		app.Version = "1.0.0"
	}

	if validate {
		if err := caddyvalidate.Validate(app.CaddyConfig); err != nil {
			return err
		}
	}

	if err := caddyManager.GetManager().AddOrUpdateApp(app.Name, &app.CaddyConfig); err != nil {
		return err
	}

	if err := SaveAppConfig(app); err != nil {
		return err
	}

	if app.BackendType == ExecBackend {
		config := &ExecBackendConfig{}

		if err := json.Unmarshal(app.BackendConfigRaw, config); err != nil {
			return err
		}

		backend.StartNewBackend(app.Name, config.Path, config.AutoReboot, strings.Split(config.Argument, " ")...)
	}

	return nil
}

func GetApps() []App {
	apps := make([]App, 0)
	infos, err := ioutil.ReadDir(config.Conf.AppConf)
	if err != nil {
		return apps
	}

	for _, info := range infos {
		if strings.HasSuffix(info.Name(), ".json") {
			if app, err := LoadAppConfig(info.Name()); err != nil {
				log.Println("Failed to load ", info.Name(), ", error: ", err.Error())
			} else {
				apps = append(apps, app)
			}
		}
	}

	return apps
}

func GetApp(name string) (App, error) {
	if CheckAppExist(name) {
		return LoadAppConfig(name + ".json")
	} else {
		return App{}, errors.New("app not found")
	}
}

func EditApp(app App, validate bool) error {
	if !CheckAppExist(app.Name) {
		return errors.New("app not found")
	}

	if app.Version == "" {
		app.Version = "1.0.0"
	}

	if validate {
		if err := caddyvalidate.Validate(app.CaddyConfig); err != nil {
			return err
		}
	}

	if err := SaveAppConfig(app); err != nil {
		return err
	}

	if err := caddyManager.GetManager().AddOrUpdateApp(app.Name, &app.CaddyConfig); err != nil {
		return err
	}

	if app.BackendType == ExecBackend {
		b := backend.GetBackend(app.Name)
		if err := b.Stop(); err != nil {
			return err
		}

		config := &ExecBackendConfig{}

		if err := json.Unmarshal(app.BackendConfigRaw, config); err != nil {
			return err
		}

		backend.StartNewBackend(app.Name, config.Path, config.AutoReboot, strings.Split(config.Argument, " ")...)
	}

	return nil
}

func DeleteApp(name string) error {
	if CheckAppExist(name) {
		if err := caddyManager.GetManager().DeleteApp(name); err != nil {
			return err
		}

		return os.Remove(config.Conf.AppConf + "/" + name + ".json")
	} else {
		return errors.New("app not found")
	}
}
