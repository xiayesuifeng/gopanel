package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/auth"
)

type Auth struct {
}

func (a *Auth) Login(ctx *gin.Context) {
	type data struct {
		password string `json:"password" binding:"required"`
	}

	d := &data{}
	if err := ctx.Bind(d); err != nil {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	}

	ctx.JSON(200, gin.H{
		"code":  200,
		"token": token,
	})
}

func (a *Auth) GetToken(ctx *gin.Context) {
	newToken, err := auth.GenerateToken()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":  200,
			"token": newToken,
		})
	}
}
