package engine

import (
	"context"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine/entity"
)

type Image interface {
	Exists(ctx context.Context, nameOrID string) (bool, error)
	Remove(ctx context.Context, nameOrID string) error
	List(ctx context.Context, all bool, filters map[string][]string) ([]*entity.Image, error)
}
