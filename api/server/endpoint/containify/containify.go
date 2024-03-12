package containify

import (
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
	r.Use(c.middleware)
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
