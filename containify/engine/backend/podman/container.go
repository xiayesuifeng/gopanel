package podman

import (
	"context"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine/entity"
)

type container struct {
	podman *Podman
}

func (c *container) Create(ctx context.Context, container *entity.Container) (containerID string, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *container) Remove(ctx context.Context, nameOrID string) error {
	//TODO implement me
	panic("implement me")
}

func (c *container) List(ctx context.Context) ([]*entity.ListContainer, error) {
	//TODO implement me
	panic("implement me")
}

func (c *container) Stop(ctx context.Context, nameOrID string) error {
	//TODO implement me
	panic("implement me")
}

func (c *container) Start(ctx context.Context, nameOrID string) error {
	//TODO implement me
	panic("implement me")
}

func (c *container) Restart(ctx context.Context, nameOrID string) error {
	//TODO implement me
	panic("implement me")
}
