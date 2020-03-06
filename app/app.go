package app

import (
	"encoding/json"
	"errors"
	"gitlab.com/xiayesuifeng/gopanel/backend"
	"gitlab.com/xiayesuifeng/gopanel/caddy"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"io/ioutil"
	"log"
	"strings"
)

const (
	GO_TYPE = iota + 1
	JAVA_TYPE
	PHP_TYPE
)

type App struct {
	Name         string          `json:"name" binding:"required"`
	CaddyConfig  json.RawMessage `json:"caddyConfig" binding:"required"`
	Type         int             `json:"type" binding:"required"`
	Path         string          `json:"path"`
	AutoReboot   bool            `json:"autoReboot"`
	BootArgument string          `json:"bootArgument"`
}

func AddApp(app App) error {
	if CheckAppExist(app.Name) {
		return errors.New("app is exist")
	}

	if err := SaveAppConfig(app); err != nil {
		return err
	}

	if err := caddy.AddServer(app.Name, app.CaddyConfig); err != nil {
		return err
	}

	if app.Type == GO_TYPE {
		backend.StartNewBackend(app.Name, app.Path, app.AutoReboot, strings.Split(app.BootArgument, " ")...)
	}

	return nil
}

func GetApps() []App {
	apps := make([]App, 0)
	infos, err := ioutil.ReadDir(core.Conf.AppConf)
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
