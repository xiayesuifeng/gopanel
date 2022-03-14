package caddy

import (
	"gitlab.com/xiayesuifeng/gopanel/core/settingStorage"
	"strconv"
)

const (
	HTTPPortKey          = "caddy/httpPort"
	HTTPSPortKey         = "caddy/httpsPort"
	ExperimentalHttp3Key = "caddy/experimentalHttp3"
	AllowH2CKey          = "caddy/allowH2C"
)

type Configuration struct {
	General General `json:"general"`
}

type General struct {
	HTTPPort  int `json:"HTTPPort"`
	HTTPSPort int `json:"HTTPSPort"`

	ExperimentalHttp3 bool `json:"experimentalHttp3"`
	AllowH2C          bool `json:"allowH2C"`
}

func InitDefaultPortConf(httpPort, httpsPort int) error {
	storage := settingStorage.GetStorage()

	if !storage.Has(HTTPPortKey) {
		if err := storage.Set(HTTPPortKey, []byte(strconv.Itoa(httpPort))); err != nil {
			return err
		}
	}

	if !storage.Has(HTTPSPortKey) {
		if err := storage.Set(HTTPSPortKey, []byte(strconv.Itoa(httpsPort))); err != nil {
			return err
		}
	}

	return nil
}

func GetConfiguration() *Configuration {
	storage := settingStorage.GetStorage()

	caddy := &Configuration{
		General: General{},
	}

	httpPort, err := strconv.Atoi(string(storage.Get(HTTPPortKey, []byte("80"))))
	if err == nil {
		caddy.General.HTTPPort = httpPort
	}

	httpsPort, err := strconv.Atoi(string(storage.Get(HTTPSPortKey, []byte("443"))))
	if err == nil {
		caddy.General.HTTPSPort = httpsPort
	}

	experimentalHttp3, err := strconv.ParseBool(string(storage.Get(ExperimentalHttp3Key, []byte("false"))))
	if err == nil {
		caddy.General.ExperimentalHttp3 = experimentalHttp3
	}

	allowH2C, err := strconv.ParseBool(string(storage.Get(AllowH2CKey, []byte("false"))))
	if err == nil {
		caddy.General.ExperimentalHttp3 = allowH2C
	}

	return caddy
}
