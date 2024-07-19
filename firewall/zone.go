package firewall

import "gitlab.com/xiayesuifeng/go-firewalld"

type ZoneStrategy string

const (
	AcceptZoneStrategy   = "ACCEPT"
	DefaultZoneStrategy  = "default"
	RejectZoneStrategy   = "REJECT"
	ContinueZoneStrategy = "CONTINUE"
	DropZoneStrategy     = "DROP"
)

type Zone struct {
	Name               string       `json:"name"`
	Description        string       `json:"description"`
	Target             ZoneStrategy `json:"target"`
	IngressPriority    int          `json:"ingressPriority"`
	EgressPriority     int          `json:"egressPriority"`
	ICMPBlocks         []string     `json:"icmpBlocks"`
	ICMPBlockInversion bool         `json:"icmpBlockInversion"`
	Masquerade         bool         `json:"masquerade"`
	Forward            bool         `json:"forward"`
	Interfaces         []string     `json:"interfaces"`
	Protocols          []string     `json:"protocols"`
}

func GetZone(name string, permanent bool) (*Zone, error) {
	conn, err := firewalld.New()
	if err != nil {
		return nil, err
	}

	if permanent {
		conn = conn.Permanent()
	}

	zone, err := conn.GetZoneByName(name)
	if err != nil {
		return nil, err
	}

	return &Zone{
		Name:               zone.Name,
		Description:        zone.Description,
		Target:             ZoneStrategy(zone.Target),
		IngressPriority:    zone.IngressPriority,
		EgressPriority:     zone.EgressPriority,
		ICMPBlocks:         zone.ICMPBlocks,
		ICMPBlockInversion: zone.ICMPBlockInversion,
		Masquerade:         zone.Masquerade,
		Forward:            zone.Forward,
		Interfaces:         zone.Interfaces,
		Protocols:          zone.Protocols,
	}, nil
}
