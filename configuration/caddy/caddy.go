package caddy

import (
	"encoding/json"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	"gitlab.com/xiayesuifeng/gopanel/core/settingStorage"
	"log"
	"strconv"
)

const (
	HTTPPortKey          = "general/httpPort"
	HTTPSPortKey         = "general/httpsPort"
	ExperimentalHttp3Key = "general/experimentalHttp3"
	AllowH2CKey          = "general/allowH2C"

	TLSKey = "tls"

	module = "caddy"
)

type Configuration struct {
	General General `json:"general"`
	TLS     TLS     `json:"tls"`
}

type General struct {
	HTTPPort  int `json:"HTTPPort"`
	HTTPSPort int `json:"HTTPSPort"`

	ExperimentalHttp3 bool `json:"experimentalHttp3"`
	AllowH2C          bool `json:"allowH2C"`
}

type TLS struct {
	DNSChallenges   map[string]caddytls.DNSChallengeConfig `json:"dnsChallenges"`
	WildcardDomains []string                               `json:"wildcardDomains"`
}

func InitDefaultPortConf(httpPort, httpsPort int) error {
	storage := settingStorage.GetStorage()

	if !storage.Has(module, HTTPPortKey) {
		if err := storage.Set(module, HTTPPortKey, []byte(strconv.Itoa(httpPort))); err != nil {
			return err
		}
	}

	if !storage.Has(module, HTTPSPortKey) {
		if err := storage.Set(module, HTTPSPortKey, []byte(strconv.Itoa(httpsPort))); err != nil {
			return err
		}
	}

	return nil
}

func GetConfiguration() *Configuration {
	storage := settingStorage.GetStorage()

	caddy := &Configuration{
		General: General{},
		TLS: TLS{
			DNSChallenges:   map[string]caddytls.DNSChallengeConfig{},
			WildcardDomains: []string{},
		},
	}

	httpPort, err := strconv.Atoi(string(storage.Get(module, HTTPPortKey, []byte("80"))))
	if err == nil {
		caddy.General.HTTPPort = httpPort
	}

	httpsPort, err := strconv.Atoi(string(storage.Get(module, HTTPSPortKey, []byte("443"))))
	if err == nil {
		caddy.General.HTTPSPort = httpsPort
	}

	experimentalHttp3, err := strconv.ParseBool(string(storage.Get(module, ExperimentalHttp3Key, []byte("false"))))
	if err == nil {
		caddy.General.ExperimentalHttp3 = experimentalHttp3
	}

	allowH2C, err := strconv.ParseBool(string(storage.Get(module, AllowH2CKey, []byte("false"))))
	if err == nil {
		caddy.General.ExperimentalHttp3 = allowH2C
	}

	if tlsRaw := storage.Get(module, TLSKey, nil); tlsRaw != nil {
		if err := json.Unmarshal(tlsRaw, &caddy.TLS); err != nil {
			log.Println("[caddy configuration] unmarshal tls struct fail:", err)
		}
	}

	return caddy
}
