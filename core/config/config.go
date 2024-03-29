package config

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
)

var (
	Conf              = &Config{}
	currentConfigFile string
)

type Config struct {
	Mode     string  `json:"mode"`
	Password string  `json:"password"`
	AppConf  string  `json:"appConf"`
	Data     string  `json:"data"`
	Secret   string  `json:"secret"`
	Panel    Panel   `json:"panel"`
	Caddy    Caddy   `json:"caddy"`
	Smtp     Smtp    `json:"smtp"`
	Netdata  Netdata `json:"netdata"`
	Log      Log     `json:"log"`
}

type Log struct {
	Level  string `json:"level"`
	Output string `json:"output"`
	Format string `json:"format"`
}

type Panel struct {
	Domain     string `json:"domain,omitempty"`
	Port       int    `json:"port,omitempty"`
	DisableSSL bool   `json:"disableSSL,omitempty"`
}

type Caddy struct {
	AdminAddress NetAddress `json:"adminAddress"`
	Conf         string     `json:"conf"`
	Data         string     `json:"data"`

	DefaultHTTPPort  int `json:"-"`
	DefaultHTTPSPort int `json:"-"`
}

type Smtp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

type Netdata struct {
	Enable bool   `json:"enable"`
	Host   string `json:"host"`
	Path   string `json:"path,omitempty"`
	SSL    bool   `json:"ssl"`
}

func ParseConf(config string) (firstLaunch bool, err error) {
	viper.SetConfigType("json")
	viper.SetConfigFile(config)

	viper.SetEnvPrefix("GOPANEL")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok || os.IsNotExist(err) {
			// Config file not found, first launch need to install
			firstLaunch = true
		} else {
			return false, err
		}
	}

	viper.SetDefault("mode", "debug")
	viper.SetDefault("data", "data")
	viper.SetDefault("appConf", "app.conf.d")
	viper.SetDefault("caddy", Caddy{
		DefaultHTTPPort:  80,
		DefaultHTTPSPort: 443,
	})
	viper.SetDefault("log", Log{
		Level:  "info",
		Output: "stderr",
		Format: "text",
	})

	err = viper.Unmarshal(Conf, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))

	currentConfigFile = config

	return
}

func SaveConf() error {
	conf, err := os.OpenFile(currentConfigFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(conf).Encode(&Conf)
}
