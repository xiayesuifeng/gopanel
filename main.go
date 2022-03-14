package main

import (
	"flag"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/app"
	"gitlab.com/xiayesuifeng/gopanel/auth"
	"gitlab.com/xiayesuifeng/gopanel/caddy"
	"gitlab.com/xiayesuifeng/gopanel/controller"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"gitlab.com/xiayesuifeng/gopanel/core/storage"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	config = flag.String("c", "config.json", "config file path")
	port   = flag.Int("p", 8080, "port")
	help   = flag.Bool("h", false, "help")
)

func main() {
	router := gin.Default()

	apiRouter := router.Group("/api", auth.AuthMiddleware)

	authRouter := apiRouter.Group("/auth")
	{
		authC := &controller.Auth{}
		authRouter.GET("/token", authC.GetToken)
		authRouter.POST("/login", authC.Login)
	}

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

	caddyRouter := apiRouter.Group("/caddy")
	{
		caddyC := &controller.Caddy{}
		caddyRouter.GET("/configuration", caddyC.GetConfiguration)
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

	err := core.ParseConf(*config)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("please config config.json")
			os.Exit(0)
		}
		log.Panicln(err)
	}

	appConf := os.Getenv("GOPANEL_APP_CONF_PATH")
	if appConf != "" {
		core.Conf.AppConf = appConf
	}
	if _, err := os.Stat(core.Conf.AppConf); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(core.Conf.AppConf, 0755)
		} else {
			log.Panicln("app.conf.d dir create failure")
		}
	}

	if err := storage.InitBaseStorage(core.Conf.Data); err != nil {
		log.Fatalln(err)
	}

	if err := caddyManager.InitManager(core.Conf.Caddy.AdminAddress); err != nil {
		log.Fatalln(err)
	}

	if err := caddy.LoadPanelConfig(strconv.Itoa(*port)); err != nil {
		log.Fatalln(err)
	}

	app.ReloadAppConfig()
}
