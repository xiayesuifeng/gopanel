/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/app"
	"gitlab.com/xiayesuifeng/gopanel/auth"
	"gitlab.com/xiayesuifeng/gopanel/controller"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"gitlab.com/xiayesuifeng/gopanel/core/storage"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	port int
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Run:   serverRun,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "server listen port")
}

func serverInit() error {
	appConf := os.Getenv("GOPANEL_APP_CONF_PATH")
	if appConf != "" {
		core.Conf.AppConf = appConf
	}
	if _, err := os.Stat(core.Conf.AppConf); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(core.Conf.AppConf, 0755)
		} else {
			log.Fatalln("app.conf.d dir create failure")
		}
	}

	if err := storage.InitBaseStorage(core.Conf.Data); err != nil {
		return err
	}

	if err := caddyManager.InitManager(core.Conf.Caddy.AdminAddress, strconv.Itoa(port)); err != nil {
		return err
	}

	app.ReloadAppConfig()

	return nil
}

func serverRun(cmd *cobra.Command, args []string) {
	if err := serverInit(); err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()

	apiRouter := router.Group("/api", auth.AuthMiddleware)

	serviceRouter := apiRouter.Group("/service")
	{
		serviceC := &controller.Service{}
		serviceRouter.GET("", serviceC.Get)
		serviceRouter.POST("/:name/:action", serviceC.Post)
	}

	webPath := os.Getenv("GOPANEL_WEB_PATH")
	if webPath == "" {
		webPath = "web"
	}
	router.Use(static.Serve("/", static.LocalFile(webPath, false)))
	router.NoRoute(func(c *gin.Context) {
		if !strings.Contains(c.Request.RequestURI, "/api") && !strings.Contains(c.Request.RequestURI, "/netdata") {
			path := strings.Split(c.Request.URL.Path, "/")
			if len(path) > 1 {
				c.File(webPath + "/index.html")
				return
			}
		}
	})

	if err := router.Run(":" + strconv.FormatInt(int64(port), 10)); err != nil {
		log.Fatalln(err)
	}
}
