package caddyManager

import (
	"errors"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

const httpApi = "/config/apps/http"

func (m *Manager) getAppConfig() (app *caddyhttp.App, err error) {
	app = &caddyhttp.App{}

	resp, err := m.httpClient.R().SetResult(app).Get(httpApi)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New("caddy admin api return message: " + string(resp.Body()))
	}

	return
}
