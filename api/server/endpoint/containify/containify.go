package containify

import (
	"encoding/json"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/containify"
	"net/http"
)

type Containify struct {
	service *containify.Containify
}

func (c *Containify) Name() string {
	return "containify"
}

func (c *Containify) Run(r router.Router) {
	r.GET("", c.Get)
	r.PUT("", c.Put)

	r.Use(c.middleware)

	r.GET("/images", c.GetImages)
	r.DELETE("/image/:nameOrID", c.RemoveImage)

	r.GET("/containers", c.GetContainers)
	r.POST("/container/:nameOrID/start", c.StartContainer)
	r.POST("/container/:nameOrID/restart", c.RestartContainer)
	r.POST("/container/:nameOrID/stop", c.StopContainer)
}

type configuration struct {
	Enabled                bool            `json:"enabled"`
	ContainerEngine        string          `json:"containerEngine"`
	ContainerEngineSetting json.RawMessage `json:"containerEngineSetting" json:"containerEngineSetting" binding:"required"`
}

func (c *Containify) Get(ctx *router.Context) error {
	engine, setting := containify.GetContainerEngine()

	return ctx.JSON(&configuration{
		Enabled:                containify.IsEnabled(),
		ContainerEngine:        engine,
		ContainerEngineSetting: setting,
	})
}

func (c *Containify) Put(ctx *router.Context) error {
	data := &configuration{}

	if err := ctx.ShouldBind(data); err != nil {
		return err
	}

	if err := containify.SetEnabled(data.Enabled); err != nil {
		return err
	}

	if err := containify.SetContainerEngine(data.ContainerEngine, data.ContainerEngineSetting); err != nil {
		return err
	}

	return ctx.NoContent()
}

func (c *Containify) middleware(ctx *router.Context) error {
	if c.service == nil {
		instance, err := containify.New()
		if err != nil {
			ctx.Abort()
			return ctx.Error(http.StatusServiceUnavailable, err.Error())
		}

		c.service = instance
	}

	ctx.Next()

	return nil
}
