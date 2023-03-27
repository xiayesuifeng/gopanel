package middleware

import (
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/core/control"
	"net/http"
	"strings"
)

func InstallMiddleware(ctx *router.Context) error {
	if strings.HasSuffix(ctx.Request.RequestURI, "/api/install") {
		ctx.Next()
		return nil
	}

	if control.Control.IsFirstLaunch() {
		ctx.Abort()
		return ctx.Error(http.StatusServiceUnavailable, "system has not been installed and initialized")
	}

	ctx.Next()
	return nil
}
