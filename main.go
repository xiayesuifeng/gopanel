package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/app"
	"gitlab.com/xiayesuifeng/gopanel/controller"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"log"
	"os"
	"strconv"
)

var (
	port = flag.Int("p", 8080, "port")
	help = flag.Bool("h", false, "help")
)

func main() {
	router := gin.Default()

	apiRouter := router.Group("/api")

	appRouter := apiRouter.Group("/app")
	{
		appC := &controller.App{}
		appRouter.GET("", appC.Gets)
		appRouter.GET("/:name", appC.Get)
		appRouter.POST("", appC.Post)
		appRouter.PUT("/:name", appC.Put)
		appRouter.DELETE("/:name", appC.Delete)
	}

	backendRouter := apiRouter.Group("/backend")
	{
		backendC := &controller.Backend{}
		backendRouter.GET("/:name", backendC.Get)
		backendRouter.GET("/:name/ws", backendC.GetWS)
	}

	if err := router.Run(":" + strconv.FormatInt(int64(*port), 10)); err != nil {
		log.Panicln(err)
	}
}

func init() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	err := core.ParseConf("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("please config config.json")
			os.Exit(0)
		}
		log.Panicln(err)
	}

	if _, err := os.Stat(core.Conf.AppConf); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(core.Conf.AppConf, 0755)
		} else {
			log.Panicln("app.conf.d dir create failure")
		}
	}

	app.ReloadAppConfig()
}
