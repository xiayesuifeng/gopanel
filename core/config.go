package core

import (
	"encoding/json"
	"os"
)

var Conf *Config

type Config struct {
	Mode     string   `json:"mode"`
	Password string   `json:"password"`
	AppConf  string   `json:"appConf"`
	Secret   string   `json:"secret"`
	Panel    Panel    `json:"panel"`
	Caddy    Caddy    `json:"caddy"`
	Db       Database `json:"database"`
	Smtp     Smtp     `json:"smtp"`
	Netdata  Netdata  `json:"netdata"`
}

type Panel struct {
	Domain         string `json:"domain,omitempty"`
	Port           int    `json:"port,omitempty"`
	AutomaticHttps bool   `json:"automaticHttps,omitempty"`
}

type Caddy struct {
	AdminAddress string `json:"adminAddress"`
	Conf         string `json:"conf"`
	Data         string `json:"data"`
}

type Database struct {
	Driver   string `json:"driver"`
	Address  string `json:"address" form:"address"`
	Port     string `json:"port" form:"port"`
	Dbname   string `json:"dbname" form:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
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
}

func ParseConf(config string) error {
	var c Config
	c.Panel.AutomaticHttps = true

	conf, err := os.Open(config)
	if err != nil {
		return err
	}
	err = json.NewDecoder(conf).Decode(&c)

	Conf = &c
	return err
}

func SaveConf() error {
	conf, err := os.OpenFile("config.json", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(conf).Encode(&Conf)
}
