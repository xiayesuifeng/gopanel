package caddymodule

import "testing"

func TestGetOfficialPluginList(t *testing.T) {
	if list, err := GetOfficialPluginList(); err != nil {
		t.Error(err)
	} else {
		for _, info := range list {
			t.Log(info)
		}
	}
}
