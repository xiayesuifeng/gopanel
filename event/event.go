package event

import (
	"slices"
	"sync"
)

type Type string

const (
	CreatedType Type = "created"
	UpdatedType Type = "updated"
	DeletedType Type = "deleted"
)

type Event struct {
	Topic   string      `json:"topic"`
	Type    Type        `json:"type"`
	Payload interface{} `json:"payload"`
}

var (
	mutex       sync.RWMutex
	subscribers = make(map[string][]chan Event)
)

func Publish(event Event) {
	mutex.RLock()
	defer mutex.RUnlock()

	if channels, ok := subscribers[event.Topic]; ok {
		for _, c := range channels {
			c <- event
		}
	}
}

func Subscribe(topic string, channel chan Event) {
	mutex.Lock()
	defer mutex.Unlock()

	if channels, ok := subscribers[topic]; ok {
		subscribers[topic] = append(channels, channel)
	} else {
		subscribers[topic] = []chan Event{channel}
	}
}

func Unsubscribe(topic string, channel chan Event) {
	mutex.Lock()
	defer mutex.Unlock()

	if channels, ok := subscribers[topic]; ok {
		subscribers[topic] = slices.DeleteFunc(channels, func(c chan Event) bool {
			return c == channel
		})
	}
}
