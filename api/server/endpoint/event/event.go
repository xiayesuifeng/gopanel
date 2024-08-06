package event

import (
	"errors"
	"github.com/gorilla/websocket"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"gitlab.com/xiayesuifeng/gopanel/event"
	"net/http"
	"slices"
)

type Event struct{}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (e *Event) Name() string {
	return "event"
}

func (e *Event) Run(r router.Router) {
	r.GET("", e.Websocket)
}

type Type string

const (
	SubscribeType   Type = "subscribe"
	UnsubscribeType Type = "unsubscribe"
)

type request struct {
	Type  Type   `json:"type"`
	Topic string `json:"topic"`
}

func (e *Event) Websocket(ctx *router.Context) error {
	headers := http.Header{}
	headers.Add("Sec-WebSocket-Protocol", ctx.GetHeader("Sec-WebSocket-Protocol"))

	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, headers)
	if err != nil {
		return err
	}
	defer conn.Close()

	exit := make(chan struct{}, 1)

	eventChan := make(chan event.Event, 10)
	topics := make([]string, 0)

	go func() {
		req := &request{}
		for {
			err := conn.ReadJSON(req)
			if err != nil {
				var closeError *websocket.CloseError
				if errors.As(err, &closeError) {
					exit <- struct{}{}
					break
				}
			}

			switch req.Type {
			case SubscribeType:
				if !slices.Contains(topics, req.Topic) {
					event.Subscribe(req.Topic, eventChan)
					topics = append(topics, req.Topic)
				}
			case UnsubscribeType:
				event.Unsubscribe(req.Topic, eventChan)
				topics = slices.DeleteFunc(topics, func(s string) bool {
					return s == req.Topic
				})
			}
		}
	}()

	for {
		select {
		case e := <-eventChan:
			err := conn.WriteJSON(e)
			if err != nil {
				return err
			}
		case <-exit:
			return nil
		}
	}
}
