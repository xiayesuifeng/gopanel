package app

import (
	"encoding/json"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"os"
)

func SaveAppConfig(app App) error {
	conf, err := os.OpenFile(core.Conf.AppConf+"/"+app.Name+".json", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(conf).Encode(&app)
}

func CheckAppExist(name string) bool {
	if _, err := os.Stat(core.Conf.AppConf + "/" + name + ".json"); err != nil {
		return !os.IsNotExist(err)
	}

	return true
}
