package caddyvalidate

import (
	"encoding/json"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"gitlab.com/xiayesuifeng/gopanel/experiments/caddyManager"
)

func Validate(config caddyManager.APPConfig) error {
	bytes, err := json.Marshal(config.Routes)
	if err != nil {
		return err
	}

	bytes = caddy.RemoveMetaFields(bytes)

	var routes []caddyhttp.Route
	err = caddy.StrictUnmarshalJSON(bytes, &routes)
	if err != nil {
		return err
	}

	return nil
}
