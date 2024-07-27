package firewall

import (
	"cmp"
	"errors"
	"slices"
	"sync"

	"gitlab.com/xiayesuifeng/go-firewalld"
)

type ZoneStrategy string

const (
	AcceptZoneStrategy   = "ACCEPT"
	DefaultZoneStrategy  = "default"
	RejectZoneStrategy   = "%%REJECT%%"
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

	return toZone(zone), nil
}

func AddZone(zone *Zone) error {
	conn, err := getConn(true)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.AddZone(&firewalld.Zone{
		Name:               zone.Name,
		Description:        zone.Description,
		Target:             string(zone.Target),
		IngressPriority:    zone.IngressPriority,
		EgressPriority:     zone.EgressPriority,
		ICMPBlocks:         zone.ICMPBlocks,
		ICMPBlockInversion: zone.ICMPBlockInversion,
		Masquerade:         zone.Masquerade,
		Forward:            zone.Forward,
		Interfaces:         zone.Interfaces,
		Protocols:          zone.Protocols,
	})
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

func GetZoneNames(permanent bool) ([]string, error) {
	conn, err := getConn(permanent)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn.GetZoneNames()
}

func GetZones(permanent bool) (result []*Zone, err error) {
	conn, err := getConn(permanent)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	names, err := conn.GetZoneNames()
	if err != nil {
		return nil, err
	}

	var errs []error

	wg := sync.WaitGroup{}
	var mutex sync.Mutex
	for _, name := range names {
		wg.Add(1)

		go func() {
			defer wg.Done()
			zone, err := conn.GetZoneByName(name)
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				errs = append(errs, err)
				return
			}

			result = append(result, toZone(zone))
		}()
	}

	wg.Wait()
	err = errors.Join(errs...)

	if err == nil {
		slices.SortFunc(result, func(a, b *Zone) int {
			return cmp.Compare[string](a.Name, b.Name)
		})
	}

	return
}

func toZone(zone *firewalld.Zone) *Zone {
	icmpBlocks := make([]string, 0)
	if zone.ICMPBlocks != nil {
		icmpBlocks = zone.ICMPBlocks
	}
	interfaces := make([]string, 0)
	if zone.Interfaces != nil {
		interfaces = zone.Interfaces
	}
	protocols := make([]string, 0)
	if zone.Protocols != nil {
		protocols = zone.Protocols
	}

	return &Zone{
		Name:               zone.Name,
		Description:        zone.Description,
		Target:             ZoneStrategy(zone.Target),
		IngressPriority:    zone.IngressPriority,
		EgressPriority:     zone.EgressPriority,
		ICMPBlocks:         icmpBlocks,
		ICMPBlockInversion: zone.ICMPBlockInversion,
		Masquerade:         zone.Masquerade,
		Forward:            zone.Forward,
		Interfaces:         interfaces,
		Protocols:          protocols,
	}
}
