package auth

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func AuthMiddleware(ctx *gin.Context) {
	if strings.HasSuffix(ctx.Request.RequestURI, "/api/auth/login") {
		ctx.Next()
		return
	}

	token := ctx.GetHeader("Authorization")

	if strings.HasPrefix(ctx.Request.RequestURI, "/api/backend") && strings.HasSuffix(ctx.Request.RequestURI, "/ws") {
		token = ctx.GetHeader("Sec-WebSocket-Protocol")
	}

	claims, err := ParseToken(token)
	if err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{
			"code":    401,
			"message": "no authorization",
		})
		ctx.Abort()
		return
	}

	if err := claims.Valid(); err != nil {
		ctx.JSON(200, gin.H{
			"code":    401,
			"message": "no authorization",
		})
		ctx.Abort()
		return
	}

	ctx.Next()
}
