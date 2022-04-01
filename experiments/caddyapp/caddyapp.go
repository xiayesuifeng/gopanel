package caddyapp

type CaddyApp interface {
	AppInfo() AppInfo
	LoadConfig(ctx Context) any
}

type AppInfo struct {
	Name string
	New  func(ctx Context) CaddyApp
}

type Context struct {
	Change chan bool
}

func (c *Context) NotifyChange() {
	c.Change <- true
}
