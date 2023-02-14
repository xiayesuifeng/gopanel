package config

import (
	"encoding/json"
	"fmt"
	"strings"
)

type NetAddress struct {
	Network string
	Address string
}

func (n *NetAddress) UnmarshalText(text []byte) error {
	return n.UnmarshalJSON(append([]byte{'"'}, append(text, '"')...))
}

func (n *NetAddress) UnmarshalJSON(bytes []byte) error {
	address := ""

	if err := json.Unmarshal(bytes, &address); err != nil {
		return err
	}

	if idx := strings.Index(address, "/"); idx >= 0 {
		n.Network = strings.ToLower(strings.TrimSpace(address[:idx]))
		if n.IsUnixNetwork() {
			n.Address = address[idx+1:]
		} else {
			n.Network = ""
			n.Address = address
		}
	}

	if n.Network == "" {
		n.Network = "tcp"
	}

	return nil
}

func (n *NetAddress) MarshalJSON() ([]byte, error) {
	return []byte(n.String()), nil
}

func (n *NetAddress) IsUnixNetwork() bool {
	return n.Network == "unix" || n.Network == "unixgram" || n.Network == "unixpacket"
}

func (n *NetAddress) String() string {
	if n.Network == "tcp" {
		return n.Address
	} else {
		return fmt.Sprintf("%s/%s", n.Network, n.Address)
	}
}
