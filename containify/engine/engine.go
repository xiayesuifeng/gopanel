package engine

import (
	"errors"
	"log"
)

var (
	containerEngines = make(map[string]Engine)
)

type Engine interface {
	New(setting []byte) error

	Container() Container
	Image() Image
}

func New(engine string, setting []byte) (Engine, error) {
	containerEngine := containerEngines[engine]
	if containerEngine == nil {
		return nil, errors.New("container engine " + engine + " not found")
	}

	if err := containerEngine.New(setting); err != nil {
		return nil, err
	}

	return containerEngine, nil
}

func Register(name string, engine Engine) {
	if containerEngines[name] != nil {
		log.Panicln("containify: container engine " + name + " already registered")
	}

	containerEngines[name] = engine
}
