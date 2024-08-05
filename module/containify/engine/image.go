package engine

import (
	"context"
	"gitlab.com/xiayesuifeng/gopanel/module/containify/engine/entity"
	"io"
)

type Image interface {
	Pull(ctx context.Context, rawImage string, progressWriter io.Writer) (string, error)
	Exists(ctx context.Context, nameOrID string) (bool, error)
	Remove(ctx context.Context, nameOrID string) error
	List(ctx context.Context, all bool, filters map[string][]string) ([]*entity.Image, error)
}
