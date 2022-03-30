package caddymodule

import "testing"

func TestGetModuleList(t *testing.T) {
	list, err := GetModuleList()
	if err != nil {
		t.Error(err)
	}

	t.Log(list)
}
