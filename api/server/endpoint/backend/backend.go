package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/backend"
	"log"
	"net/http"
)

type Backend struct {
}

func (b *Backend) Name() string {
	return "backend"
}

func (b *Backend) Run(r router.Router) {
	r.GET("/:name", b.Get)
	r.GET("/:name/ws", b.GetWS)
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (b *Backend) Get(ctx *router.Context) error {
	{
		b := backend.GetBackend(ctx.Param("name"))
		return ctx.JSON(gin.H{
			"status": b.Status,
			"log":    b.Log.String(),
		})
	}
}

func (b *Backend) GetWS(ctx *router.Context) error {
	headers := http.Header{}
	headers.Add("Sec-WebSocket-Protocol", ctx.GetHeader("Sec-WebSocket-Protocol"))

	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, headers)
	if err != nil {
		log.Println("Failed to set websocket upgrade: %+v", err)
		return nil
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
	return nil
}
