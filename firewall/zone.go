package firewall

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
	conn, err := getConn(permanent)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

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

// UpdateZone update zone setting, target and description field only change in permanent
func UpdateZone(zone *Zone, permanent bool) error {
	conn, err := getConn(permanent)
	if err != nil {
		return err
	}
	defer conn.Close()

	originZone, err := conn.GetZoneByName(zone.Name)
	if err != nil {
		return err
	}

	originZone.Description = zone.Description
	originZone.Target = string(zone.Target)
	originZone.IngressPriority = zone.IngressPriority
	originZone.EgressPriority = zone.EgressPriority
	originZone.ICMPBlocks = zone.ICMPBlocks
	originZone.ICMPBlockInversion = zone.ICMPBlockInversion
	originZone.Masquerade = zone.Masquerade
	originZone.Forward = zone.Forward
	originZone.Interfaces = zone.Interfaces
	originZone.Protocols = zone.Protocols

	return conn.UpdateZone(originZone)
}
