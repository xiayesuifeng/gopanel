package middleware

import (
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/auth"
	"log"
	"strings"
)

func AuthMiddleware(ctx *router.Context) error {
	if strings.HasSuffix(ctx.Request.RequestURI, "/api/auth/login") || strings.HasSuffix(ctx.Request.RequestURI, "/api/install") {
		ctx.Next()
		return nil
	}

	token := ctx.GetHeader("Authorization")

	if ctx.GetHeader("Upgrade") == "websocket" {
		token = ctx.GetHeader("Sec-WebSocket-Protocol")
	}

	claims, err := auth.ParseToken(token)
	if err != nil {
		log.Println(err)
		ctx.Abort()
		return ctx.Error(401, "no authorization")
	}

	if err := claims.Valid(); err != nil {
		ctx.Abort()
		return ctx.Error(401, "no authorization")
	}

	ctx.Next()

	return nil
}
