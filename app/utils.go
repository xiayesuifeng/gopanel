package app

import (
	"encoding/json"
	"gitlab.com/xiayesuifeng/gopanel/backend"
	"gitlab.com/xiayesuifeng/gopanel/caddy"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func SaveAppConfig(app App) error {
	conf, err := os.OpenFile(core.Conf.AppConf+"/"+app.Name+".json", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(conf).Encode(&app)
}

func LoadAppConfig(name string) (App, error) {
	app := App{}
	conf, err := os.OpenFile(core.Conf.AppConf+"/"+name, os.O_RDONLY, 0644)
	if err != nil {
		return app, err
	}
	if err := json.NewDecoder(conf).Decode(&app); err != nil {
		return app, err
	}

	return app, nil
}

func CheckAppExist(name string) bool {
	if _, err := os.Stat(core.Conf.AppConf + "/" + name + ".json"); err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

func ReloadAppConfig() {
	files, _ := ioutil.ReadDir(core.Conf.AppConf)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			app, err := LoadAppConfig(file.Name())
			if err != nil {
				log.Println(err)
			} else {
				if caddy.CheckServerExist(app.Name) {
					err = caddy.EditServer(app.Name, app.CaddyConfig)
				} else {
					err = caddy.AddServer(app.Name, app.CaddyConfig)
				}

				if err != nil {
					log.Println(err)
				}
				backend.StartNewBackend(app.Name, app.Path, app.AutoReboot, strings.Split(app.BootArgument, " ")...)
			}
		}
	}
}
