package core

import (
	"context"
	"gitlab.com/xiayesuifeng/gopanel/api/server"
	"gitlab.com/xiayesuifeng/gopanel/app"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"gitlab.com/xiayesuifeng/gopanel/core/control"
	"gitlab.com/xiayesuifeng/gopanel/core/storage"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
	"gitlab.com/xiayesuifeng/gopanel/web"
	"golang.org/x/exp/slog"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Core struct {
	listenPort int
	server     *server.Server

	firstLaunch bool

	starting  bool
	closeLock sync.Mutex

	done chan bool
}

func New(port int) (*Core, error) {
	log.Println("[core] initialize...")

	// initialize logger
	var logOut io.Writer
	switch config.Conf.Log.Output {
	case "stderr":
		logOut = os.Stderr
	case "stdout":
		logOut = os.Stdout
	default:
		logFile, err := os.OpenFile(config.Conf.Log.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		logOut = logFile
	}

	level := slog.LevelWarn
	switch config.Conf.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	logOpts := slog.HandlerOptions{Level: level}
	logger := slog.New(logOpts.NewTextHandler(logOut))
	if config.Conf.Log.Format == "json" {
		logger = slog.New(logOpts.NewJSONHandler(logOut))
	}
	slog.SetDefault(logger)

	core := &Core{
		listenPort: port,
		server:     server.NewServer(web.Assets()),
		done:       make(chan bool),
	}

	control.Control = core

	return core, nil
}

func (c *Core) Start(ctx context.Context) error {
	slog.Info("[core] starting...")

	c.firstLaunch = ctx.Value("firstLaunch").(bool)

	if c.firstLaunch {
		slog.Info("[core] first launch, skip storage and caddy manager init")
	} else {
		appConf := os.Getenv("GOPANEL_APP_CONF_PATH")
		if appConf != "" {
			config.Conf.AppConf = appConf
		}
		if _, err := os.Stat(config.Conf.AppConf); err != nil {
			if os.IsNotExist(err) {
				os.MkdirAll(config.Conf.AppConf, 0755)
			} else {
				log.Fatalln("app.conf.d dir create failure")
			}
		}

		if err := storage.InitBaseStorage(config.Conf.Data); err != nil {
			return err
		}

		if err := caddyManager.InitManager(config.Conf.Caddy.AdminAddress, strconv.Itoa(c.listenPort)); err != nil {
			return err
		}

		app.ReloadAppConfig()
	}

	c.starting = true

	if err := c.server.Run(":" + strconv.FormatInt(int64(c.listenPort), 10)); err != nil && err != http.ErrServerClosed {
		return err
	}

	<-c.done

	return nil
}

func (c *Core) Close() error {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()
	if !c.starting {
		slog.Info("[core] core already closed")
		return nil
	}

	slog.Info("[core] closing...")

	if err := c.server.Shutdown(context.Background()); err != nil {
		return err
	}

	c.done <- true

	slog.Info("[core] close")
	c.starting = false
	return nil
}

func (c *Core) IsFirstLaunch() bool {
	return c.firstLaunch
}
