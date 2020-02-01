package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/xiayesuifeng/gopanel/backend"
	"log"
)

type Backend struct {
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (b *Backend) GetWS(ctx *gin.Context) {
	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	defer conn.Close()

	exit := make(chan bool, 1)

	go func() {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				exit <- true
			}
		}
	}()

	go func() {
		b := backend.GetBackend(ctx.Param("name"))
		for {
			err = conn.WriteJSON(gin.H{
				"status": b.Status,
				"log":    b.Log.String(),
			})
			if err != nil {
				exit <- true
			}
			<-b.Notify
		}
	}()

	<-exit
}
