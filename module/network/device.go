package network

import (
	"net"
)

type Device struct {
	Index int       `json:"index"`
	MTU   int       `json:"mtu"`
	Name  string    `json:"name"`
	MAC   string    `json:"mac"`
	Flags net.Flags `json:"flags"`
}

func GetDevices() ([]*Device, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	devices := make([]*Device, 0, len(interfaces))
	for _, iface := range interfaces {
		devices = append(devices, &Device{
			Index: iface.Index,
			MTU:   iface.MTU,
			Name:  iface.Name,
			MAC:   iface.HardwareAddr.String(),
			Flags: iface.Flags,
		})
	}

	return devices, nil
}
