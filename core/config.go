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
	Caddy    Caddy    `json:"caddy"`
	Db       Database `json:"database"`
	Smtp     Smtp     `json:"smtp"`
}

type Caddy struct {
	AdminAddress string `json:"admin_address"`
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

func ParseConf(config string) error {
	var c Config

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
