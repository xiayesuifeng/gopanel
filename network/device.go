package network

import (
	"net"
)

type Device struct {
	net.Interface
}

func GetDevices() ([]*Device, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	devices := make([]*Device, 0, len(interfaces))
	for _, iface := range interfaces {
		devices = append(devices, &Device{iface})
	}

	return devices, nil
}
