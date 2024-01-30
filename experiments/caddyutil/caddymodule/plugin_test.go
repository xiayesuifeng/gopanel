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

func TestInstallPlugin(t *testing.T) {
	if err := InstallPlugin("github.com/caddy-dns/cloudflare"); err != nil {
		t.Error(err)
	} else {
		t.Log("install complete")
	}
}

func TestRemovePlugin(t *testing.T) {
	if err := RemovePlugin("github.com/caddy-dns/cloudflare"); err != nil {
		t.Error(err)
	} else {
		t.Log("remove complete")
	}
}
