package install

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/auth"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"gitlab.com/xiayesuifeng/gopanel/core/control"
	"golang.org/x/exp/slog"
	"net"
)

type Install struct {
}

func (i *Install) Name() string {
	return "install"
}

func (i *Install) Run(r router.Router) {
	r.POST("", i.Install)
}

type Request struct {
	Password string `json:"password" binding:"required"`
	Caddy    struct {
		HTTPPort     int               `json:"httpPort" binding:"required,gte=1,lte=65535"`
		HTTPSPort    int               `json:"httpsPort" binding:"required,gte=1,lte=65535"`
		AdminAddress config.NetAddress `json:"adminAddress" binding:"required"`
		ConfPath     string            `json:"confPath" binding:"required"`
		DataPath     string            `json:"dataPath" binding:"required"`
	} `json:"caddy"`
	Panel struct {
		Domain     string `json:"domain" binding:"required_without=Port"`
		Port       int    `json:"port" binding:"required_without=Domain"`
		DisableSSL bool   `json:"disableSSL"`
	} `json:"panel"`
}

func (i *Install) Install(ctx *router.Context) error {
	req := &Request{}
	if err := ctx.ShouldBind(req); err != nil {
		return ctx.Error(400, err.Error())
	}

	config.Conf.Password = auth.EncryptionPassword(req.Password)
	config.Conf.Caddy.AdminAddress = req.Caddy.AdminAddress
	config.Conf.Caddy.Conf = req.Caddy.ConfPath
	config.Conf.Caddy.Data = req.Caddy.DataPath
	config.Conf.Caddy.DefaultHTTPPort = req.Caddy.HTTPPort
	config.Conf.Caddy.DefaultHTTPSPort = req.Caddy.HTTPSPort
	if len(req.Panel.Domain) > 0 {
		config.Conf.Panel.Domain = req.Panel.Domain
	} else {
		config.Conf.Panel.Port = req.Panel.Port
	}
	config.Conf.Panel.DisableSSL = req.Panel.DisableSSL

	if err := config.SaveConf(); err != nil {
		return ctx.Error(500, err.Error())
	}

	redirectURL := "https://"
	port := req.Caddy.HTTPSPort
	if req.Panel.DisableSSL {
		redirectURL = "http://"
		port = req.Caddy.HTTPPort
	}

	if len(req.Panel.Domain) > 0 {
		redirectURL += req.Panel.Domain
	} else {
		host, _, err := net.SplitHostPort(ctx.Request.Host)
		if err != nil {
			return ctx.Error(500, err.Error())
		}
		redirectURL += host
	}

	if port != 80 && port != 443 {
		redirectURL += fmt.Sprintf(":%d", port)
	}

	go func() {
		err := control.Control.Close()
		if err != nil {
			slog.Error("[install] failed to close api server", err)
		}

		err = control.Control.Start(context.WithValue(ctx, "firstLaunch", false))
		if err != nil {
			slog.Error("[install] failed to start api server, rollback to first launch mode", err)
			err = control.Control.Start(context.WithValue(ctx, "firstLaunch", true))
			if err != nil {
				slog.Error("[install] failed to start api server in first launch mode", err)
				panic(err)
			}
		}
	}()

	return ctx.JSON(gin.H{"redirectURL": redirectURL})
}
