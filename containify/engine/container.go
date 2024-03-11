package engine

import (
	"context"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine/entity"
)

type Container interface {
	Create(ctx context.Context, container *entity.Container) (containerID string, err error)
	Remove(ctx context.Context, nameOrID string) error
	List(ctx context.Context) ([]*entity.ListContainer, error)
	Stop(ctx context.Context, nameOrID string) error
	Start(ctx context.Context, nameOrID string) error
	Restart(ctx context.Context, nameOrID string) error
}
