package caddyManager

import (
	"errors"
)

func (m Manager) postCaddyObject(uri string, object interface{}) error {
	resp, err := m.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(object).
		Post(uri)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New(string(resp.Body()))
	} else {
		return nil
	}
}
