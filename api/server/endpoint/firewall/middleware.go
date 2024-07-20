package firewall

import "gitlab.com/xiayesuifeng/gopanel/api/server/router"

func permanentMiddleware(ctx *router.Context) error {
	ctx.Set("permanent", ctx.Query("permanent") == "true")

	ctx.Next()
	return nil
}

func permanent(ctx *router.Context) bool {
	return ctx.GetBool("permanent")
}
