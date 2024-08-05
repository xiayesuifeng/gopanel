package app

import (
	"encoding/json"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
	"gitlab.com/xiayesuifeng/gopanel/module/backend"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func SaveAppConfig(app App) error {
	conf, err := os.OpenFile(config.Conf.AppConf+"/"+app.Name+".json", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(conf).Encode(&app)
}

func LoadAppConfig(name string) (App, error) {
	app := App{}
	conf, err := os.OpenFile(config.Conf.AppConf+"/"+name, os.O_RDONLY, 0644)
	if err != nil {
		return app, err
	}
	if err := json.NewDecoder(conf).Decode(&app); err != nil {
		return app, err
	}

	// Convert old version app format to new format
	if app.Version == "" {
		app, err = adaptOldAppToNewVersion(app)
		if err != nil {
			return App{}, err
		}

		if err = SaveAppConfig(app); err != nil {
			return App{}, err
		}
	}

	return app, nil
}

// adaptOldAppToNewVersion Adapt old version App to 1.0.0 format
func adaptOldAppToNewVersion(app App) (App, error) {
	hosts := caddyhttp.MatchHost{}
	routes := caddyhttp.RouteList{}

	for _, route := range app.CaddyConfig.Routes {
		for _, moduleMap := range route.MatcherSetsRaw {
			oldHosts := caddyhttp.MatchHost{}

			if hostsRaw, ok := moduleMap["host"]; ok {
				if err := json.Unmarshal(hostsRaw, &oldHosts); err != nil {
					return App{}, err
				}

				hosts = append(hosts, oldHosts...)
			}
		}

		for _, handle := range route.HandlersRaw {
			oldHandler := caddyManager.SubrouteHandle{}

			if err := json.Unmarshal(handle, &oldHandler); err != nil {
				return App{}, err
			}

			if oldHandler.Handler == "subroute" {
				routes = append(routes, oldHandler.Routes...)
			}
		}
	}

	if app.Path != "" {
		app.BackendType = ExecBackend

		bConfig := ExecBackendConfig{
			AutoReboot: app.AutoReboot,
			Path:       app.Path,
			Argument:   app.BootArgument,
		}

		bConfigRaw, err := json.Marshal(&bConfig)
		if err != nil {
			return App{}, err
		}

		app.BackendConfigRaw = bConfigRaw

		// reset deprecated field to zero value
		app.AutoReboot = false
		app.Path = ""
		app.BootArgument = ""
	} else {
		app.BackendType = NoneBackend
	}

	app.Version = "1.0.0"
	app.CaddyConfig.Domain = hosts
	app.CaddyConfig.Routes = routes

	return app, nil
}

func CheckAppExist(name string) bool {
	if _, err := os.Stat(config.Conf.AppConf + "/" + name + ".json"); err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

func ReloadAppConfig() {
	files, _ := ioutil.ReadDir(config.Conf.AppConf)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			app, err := LoadAppConfig(file.Name())
			if err != nil {
				log.Println(err)
			} else {
				manager := caddyManager.GetManager()

				err := manager.AddOrUpdateApp(app.Name, &app.CaddyConfig)
				if err != nil {
					log.Println(err)
				}

				if app.BackendType == ExecBackend {
					config := &ExecBackendConfig{}

					if err := json.Unmarshal(app.BackendConfigRaw, config); err != nil {
						log.Println(err)
					} else {
						backend.StartNewBackend(app.Name, config.Path, config.AutoReboot, strings.Split(config.Argument, " ")...)

					}
				}
			}
		}
	}
}
