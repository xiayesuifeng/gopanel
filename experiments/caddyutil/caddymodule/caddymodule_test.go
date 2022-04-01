package caddymodule

import "testing"

func TestGetModuleList(t *testing.T) {
	list, err := GetModuleList()
	if err != nil {
		t.Error(err)
	}

	t.Log(list)
}

func TestModuleList_HasNonStandardModule(t *testing.T) {
	list, err := GetModuleList()
	if err != nil {
		t.Error(err)
	}

	if list.HasNonStandardModule("dns.providers.cloudflare") {
		t.Log("dns.providers.cloudflare found")
	} else {
		t.Log("dns.providers.cloudflare not found")
	}
}

func TestModuleList_HasPackage(t *testing.T) {
	list, err := GetModuleList()
	if err != nil {
		t.Error(err)
	}

	if list.HasPackage("github.com/caddy-dns/cloudflare") {
		t.Log("github.com/caddy-dns/cloudflare found")
	} else {
		t.Log("github.com/caddy-dns/cloudflare not found")
	}
}
