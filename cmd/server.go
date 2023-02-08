/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gitlab.com/xiayesuifeng/gopanel/app"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"gitlab.com/xiayesuifeng/gopanel/core/storage"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
	"log"
	"os"
	"strconv"
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

	if err := router.Run(":" + strconv.FormatInt(int64(port), 10)); err != nil {
		log.Fatalln(err)
	}
}
