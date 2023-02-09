package core

import (
	"context"
)

type Core struct {
}

func New() (*Core, error) {
	return &Core{}, nil
}

func (c *Core) Start(ctx context.Context) error {
	return nil
}

func (c *Core) Close() error {
	return nil
}
