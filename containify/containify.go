package containify

import (
	"errors"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine"
)

type Containify struct {
	engine engine.Engine
}

func (c *Containify) ContainerEngine() engine.Engine {
	return c.engine
}

func New() (*Containify, error) {
	if !IsEnabled() {
		return nil, errors.New("containify is not enabled")
	}

	engineName, setting := GetContainerEngine()
	if engineName == "" {
		return nil, errors.New("container engine is not set")
	}

	e, err := engine.New(engineName, setting)
	if err != nil {
		return nil, err
	}

	instance := &Containify{engine: e}

	return instance, nil
}
