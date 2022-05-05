package caddyapp

import "gitlab.com/xiayesuifeng/gopanel/experiments/caddyutil/caddymodule"

type CaddyApp interface {
	AppInfo() AppInfo
	LoadConfig(ctx Context) any
}

type AppInfo struct {
	Name string
	New  func(ctx Context) CaddyApp
}

type Context struct {
	Change     chan bool
	ModuleList *caddymodule.ModuleList
}

func (c *Context) NotifyChange() {
	c.Change <- true
}
