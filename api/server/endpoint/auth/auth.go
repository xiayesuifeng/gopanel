package auth

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/auth"
	"gitlab.com/xiayesuifeng/gopanel/core"
)

type Auth struct {
}

func (a *Auth) Name() string {
	return "auth"
}

func (a *Auth) Run(r router.Router) {
	r.GET("/token", a.GetToken)
	r.POST("/login", a.Login)
}

func (a *Auth) Login(ctx *router.Context) error {
	type data struct {
		Password string `json:"password" binding:"required"`
	}

	d := &data{}
	if err := ctx.ShouldBind(d); err != nil {
		return ctx.Error(400, err.Error())
	}

	if core.EncryptionPassword(d.Password) != core.Conf.Password {
		return ctx.Error(400, "password error")
	}

	token, err := auth.GenerateToken()
	if err != nil {
		return err
	}

	return ctx.JSON(gin.H{"token": token})
}

func (a *Auth) GetToken(ctx *router.Context) error {
	newToken, err := auth.GenerateToken()
	if err != nil {
		return err
	}

	return ctx.JSON(gin.H{"token": newToken})
}
