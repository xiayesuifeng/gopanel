package control

import "context"

type control interface {
	IsFirstLaunch() bool
	Start(ctx context.Context) error
	Close() error
}

var Control control
