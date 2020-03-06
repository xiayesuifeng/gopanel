package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/auth"
	"gitlab.com/xiayesuifeng/gopanel/core"
)

type Auth struct {
}

func (a *Auth) Login(ctx *gin.Context) {
	type data struct {
		Password string `json:"password" binding:"required"`
	}

	d := &data{}
	if err := ctx.ShouldBind(d); err != nil {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	if core.EncryptionPassword(d.Password) != core.Conf.Password {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": "password error",
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
