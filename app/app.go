package app

import (
	"encoding/json"
	"errors"
	"gitlab.com/xiayesuifeng/gopanel/backend"
	"gitlab.com/xiayesuifeng/gopanel/caddy"
	"gitlab.com/xiayesuifeng/gopanel/caddy/config"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	GO_TYPE = iota + 1
	JAVA_TYPE
	PHP_TYPE
	OTHER_TYPE
)

type App struct {
	Name         string            `json:"name"`
	CaddyConfig  config.ServerType `json:"caddyConfig" binding:"required"`
	Type         int               `json:"type" binding:"required"`
	Path         string            `json:"path"`
	AutoReboot   bool              `json:"autoReboot"`
	BootArgument string            `json:"bootArgument"`
}

func AddApp(app App, validate bool) error {
	if CheckAppExist(app.Name) {
		return errors.New("app is exist")
	}

	if validate {
		bytes, err := json.Marshal(app.CaddyConfig)
		if err != nil {
			return err
		}

		if err := caddy.ValidateConfig(bytes); err != nil {
			return err
		}

		go caddy.AddAdaptServerToRoute(app.Name, app.CaddyConfig)
	} else if err := caddy.AddAdaptServerToRoute(app.Name, app.CaddyConfig); err != nil {
		return err
	}

	if err := SaveAppConfig(app); err != nil {
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

	if err := SaveAppConfig(app); err != nil {
		return err
	}

	if validate {
		go caddy.EditAdaptServerToRoute(app.Name, app.CaddyConfig)
	} else if err := caddy.EditAdaptServerToRoute(app.Name, app.CaddyConfig); err != nil {
		return err
	}

	if app.Type == GO_TYPE {
		b := backend.GetBackend(app.Name)
		if err := b.Stop(); err != nil {
			return err
		}

		backend.StartNewBackend(app.Name, app.Path, app.AutoReboot, strings.Split(app.BootArgument, " ")...)
	}

	return nil
}

func DeleteApp(name string, validate bool) error {
	if CheckAppExist(name) {
		if validate {
			go caddy.DeleteAdaptServerToRoute(name)
		} else if err := caddy.DeleteServer(name); err != nil {
			return err
		}

		return os.Remove(core.Conf.AppConf + "/" + name + ".json")
	} else {
		return errors.New("app not found")
	}
}
