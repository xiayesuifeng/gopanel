package caddyManager

import "github.com/caddyserver/caddy/v2/modules/caddyhttp"

type SubrouteHandle struct {
	Handler string              `json:"handler"`
	Routes  caddyhttp.RouteList `json:"routes"`
}

func NewSubrouteHandle(routes []caddyhttp.Route) *SubrouteHandle {
	return &SubrouteHandle{
		Handler: "subroute",
		Routes:  routes,
	}
}
