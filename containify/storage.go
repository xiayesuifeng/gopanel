package containify

import (
	"gitlab.com/xiayesuifeng/gopanel/core/settingStorage"
)

const (
	module                    = "containify"
	enabledKey                = "enabled"
	containerEngineKey        = "containerEngine"
	containerEngineSettingKey = "containerEngineSetting"
)

func IsEnabled() bool {
	return string(settingStorage.GetStorage().Get(module, enabledKey, []byte("false"))) == "true"
}

func SetEnabled(enabled bool) error {
	if enabled {
		return settingStorage.GetStorage().Set(module, enabledKey, []byte("true"))
	} else {
		return settingStorage.GetStorage().Set(module, enabledKey, []byte("false"))
	}
}

func GetContainerEngine() (engine string, setting []byte) {
	engine = string(settingStorage.GetStorage().Get(module, containerEngineKey, []byte("")))

	setting = settingStorage.GetStorage().Get(module, containerEngineSettingKey, []byte("{}"))
	return
}

func SetContainerEngine(engine string, setting []byte) error {
	if err := settingStorage.GetStorage().Set(module, containerEngineKey, []byte(engine)); err != nil {
		return err
	}

	return settingStorage.GetStorage().Set(module, containerEngineSettingKey, setting)
}
