package core

import (
	"context"
	"gitlab.com/xiayesuifeng/gopanel/app"
	"gitlab.com/xiayesuifeng/gopanel/core/storage"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
	"log"
	"os"
	"strconv"
)

type Core struct {
	listenPort int
}

func New(port int) (*Core, error) {
	return &Core{listenPort: port}, nil
}

func (c *Core) Start(ctx context.Context) error {
	appConf := os.Getenv("GOPANEL_APP_CONF_PATH")
	if appConf != "" {
		Conf.AppConf = appConf
	}
	if _, err := os.Stat(Conf.AppConf); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(Conf.AppConf, 0755)
		} else {
			log.Fatalln("app.conf.d dir create failure")
		}
	}

	if err := storage.InitBaseStorage(Conf.Data); err != nil {
		return err
	}

	if err := caddyManager.InitManager(Conf.Caddy.AdminAddress, strconv.Itoa(c.listenPort)); err != nil {
		return err
	}

	app.ReloadAppConfig()

	return nil
}

func (c *Core) Close() error {
	return nil
}
