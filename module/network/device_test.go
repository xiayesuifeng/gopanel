package network

import "testing"

func TestGetDevices(t *testing.T) {
	devices, err := GetDevices()
	if err != nil {
		t.Fatal(err)
	}

	for _, device := range devices {
		t.Log(device)
	}
}
