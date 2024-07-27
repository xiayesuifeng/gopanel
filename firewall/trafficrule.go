package firewall

import (
	"encoding/json"
	"fmt"
	"gitlab.com/xiayesuifeng/go-firewalld"
	"slices"
	"strings"
)

type RuleStrategy int

const (
	AcceptRuleStrategy RuleStrategy = iota
	RejectRuleStrategy
	DropRuleStrategy
)

type RuleType string

const (
	ServiceRuleType     RuleType = "service"
	PortRuleType        RuleType = "port"
	ProtocolRuleType    RuleType = "protocol"
	MasqueradeRuleType  RuleType = "masquerade"
	IcmpBlockRuleType   RuleType = "icmp-block"
	ForwardPortRuleType RuleType = "forward-port"
	SourcePortRuleType  RuleType = "source-port"
)

type TrafficRule struct {
	// Family ipv4 or ipv6, empty means both
	Family string `json:"family"`
	// SrcAddr source address
	SrcAddr       string `json:"srcAddr,omitempty"`
	SrcAddrInvert bool   `json:"srcAddrInvert,omitempty"`
	// DestAddr destination address
	DestAddr       string          `json:"destAddr,omitempty"`
	DestAddrInvert bool            `json:"destAddrInvert,omitempty"`
	Strategy       RuleStrategy    `json:"strategy"`
	Type           RuleType        `json:"type"`
	Value          json.RawMessage `json:"value"`
	Log            RuleLog         `json:"log"`
	Audit          bool            `json:"audit"`
}

type RuleLog struct {
	Enabled bool   `json:"enabled"`
	Prefix  string `json:"prefix"`
	// Level emerg、alert、crit、error、warning、notice、info or debug
	Level string `json:"level"`
	Limit string `json:"limit"`
}

type Port struct {
	// Port number or range (8080-8085)
	Port string `json:"port"`
	// Protocol tcp or udp
	Protocol string `json:"protocol"`
}

var ruleElement = []string{
	"destination",
	"service",
	"port",
	"protocol",
	"masquerade",
	"icmp-block",
	"forward-port",
	"source-port",
	"accept",
	"drop",
	"reject",
}

func GetTrafficRules(zone string, permanent bool) ([]*TrafficRule, error) {
	conn, err := getConn(permanent)
	if err != nil {
		return nil, err
	}

	result, err := conn.GetZoneByName(zone)
	if err != nil {
		return nil, err
	}

	rules := make([]*TrafficRule, 0)

	for _, service := range result.Services {
		bytes, err := json.Marshal(&service)
		if err != nil {
			return nil, err
		}

		rules = append(rules, &TrafficRule{
			Type:     ServiceRuleType,
			Value:    bytes,
			Strategy: AcceptRuleStrategy,
		})
	}

	for _, port := range result.Ports {
		bytes, err := json.Marshal(&Port{
			port.Port,
			port.Protocol,
		})
		if err != nil {
			return nil, err
		}

		rules = append(rules, &TrafficRule{
			Type:     PortRuleType,
			Value:    bytes,
			Strategy: AcceptRuleStrategy,
		})
	}

	for _, port := range result.SourcePorts {
		bytes, err := json.Marshal(&Port{
			port.Port,
			port.Protocol,
		})
		if err != nil {
			return nil, err
		}

		rules = append(rules, &TrafficRule{
			Type:     SourcePortRuleType,
			Value:    bytes,
			Strategy: AcceptRuleStrategy,
		})
	}

	for _, richRule := range result.RichRules {
		rule := &TrafficRule{}

		args := strings.Split(richRule, " ")
		for i := 1; i < len(args); i++ {
			vals := strings.Split(args[i], "=")

			switch vals[0] {
			case "family":
				rule.Family = strings.Trim(vals[1], "\"")
			case "source":
				i++

				for {
					if args[i] == "NOT" {
						rule.SrcAddrInvert = true
					} else if strings.HasPrefix(args[i], "address") {
						vals = strings.Split(args[i], "=")
						rule.SrcAddr = strings.Trim(vals[1], "\"")
					} else if strings.HasPrefix(args[i], "mac") {
						// TODO
					} else if strings.HasPrefix(args[i], "ipset") {
						// TODO
					}

					end := false
					for _, element := range ruleElement {
						if strings.HasSuffix(args[i+1], element) {
							end = true
						}
					}

					if end {
						break
					} else {
						i++
					}
				}
			case "destination":
				i++
				if args[i] == "NOT" {
					rule.DestAddrInvert = true
					i++
				}

				if strings.HasPrefix(args[i], "address") {
					vals = strings.Split(args[i], "=")
					rule.DestAddr = strings.Trim(vals[1], "\"")
				}
			case "service":
				i++
				vals = strings.Split(args[i], "=")
				rule.Type = ServiceRuleType
				rule.Value = []byte(vals[1])
			case "port":
				vals = strings.Split(args[i+1], "=")
				port := strings.Trim(vals[1], "\"")
				vals = strings.Split(args[i+2], "=")
				rule.Type = PortRuleType
				bytes, err := json.Marshal(&Port{
					port, strings.Trim(vals[1], "\""),
				})
				if err != nil {
					return nil, err
				}
				rule.Value = bytes

				i += 2
			case "protocol":
				i++
				vals = strings.Split(args[i], "=")
				rule.Type = ProtocolRuleType
				rule.Value = []byte(vals[1])
			case "masquerade":
				rule.Type = MasqueradeRuleType
			case "icmp-block":
				i++
				vals = strings.Split(args[i], "=")
				rule.Type = IcmpBlockRuleType
				rule.Value = []byte(vals[1])
			case "forward-port":
				rule.Type = ForwardPortRuleType
				// TODO
			case "source-port":
				vals = strings.Split(args[i+1], "=")
				port := strings.Trim(vals[1], "\"")
				vals = strings.Split(args[i+2], "=")
				rule.Type = SourcePortRuleType
				bytes, err := json.Marshal(&Port{
					port, strings.Trim(vals[1], "\""),
				})
				if err != nil {
					return nil, err
				}
				rule.Value = bytes

				i += 2
			case "log":
				rule.Log.Enabled = true
				i++
				if strings.HasPrefix(args[i], "prefix") {
					rule.Log.Prefix = strings.Trim(strings.Split(args[i], "=")[1], "\"")
					i++
				}
				if strings.HasPrefix(args[i], "level") {
					rule.Log.Level = strings.Trim(strings.Split(args[i], "=")[1], "\"")
					i++
				}
				if strings.HasPrefix(args[i], "limit") {
					i++
					rule.Log.Limit = strings.Trim(strings.Split(args[i], "=")[1], "\"")
				}
			case "audit":
				rule.Audit = true
			case "accept":
				rule.Strategy = AcceptRuleStrategy
			case "reject":
				rule.Strategy = RejectRuleStrategy
			case "drop":
				rule.Strategy = DropRuleStrategy
			}
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func AddTrafficRule(zone string, rule *TrafficRule, permanent bool) error {
	if rule == nil {
		return nil
	}

	if rule.Type == ForwardPortRuleType {
		// ForwardPortRuleType should be moved to addPortForwarding func handle
		return nil
	}

	conn, err := getConn(permanent)
	if err != nil {
		return err
	}

	setting, err := conn.GetZoneByName(zone)
	if err != nil {
		return err
	}

	if rule.needRichRule() {
		richRule, err := rule.toRichRule()
		if err != nil {
			return err
		}

		setting.RichRules = append(setting.RichRules, richRule)
	} else {
		switch rule.Type {
		case ServiceRuleType:
			var val string
			err := json.Unmarshal(rule.Value, &val)
			if err != nil {
				return err
			}
			setting.Services = append(setting.Services, val)
		case PortRuleType:
			var port Port
			err := json.Unmarshal(rule.Value, &port)
			if err != nil {
				return err
			}

			setting.Ports = append(setting.Ports, &firewalld.Port{
				Port:     port.Port,
				Protocol: port.Protocol,
			})
		case ProtocolRuleType:
			var val string
			err := json.Unmarshal(rule.Value, &val)
			if err != nil {
				return err
			}

			setting.Protocols = append(setting.Protocols, val)
		case MasqueradeRuleType:
			setting.Masquerade = true
		case IcmpBlockRuleType:
			var val string
			err := json.Unmarshal(rule.Value, &val)
			if err != nil {
				return err
			}

			setting.ICMPBlocks = append(setting.Protocols, val)
		case SourcePortRuleType:
			var port Port
			err := json.Unmarshal(rule.Value, &port)
			if err != nil {
				return err
			}
			setting.SourcePorts = append(setting.SourcePorts, &firewalld.Port{
				Port:     port.Port,
				Protocol: port.Protocol,
			})
		}
	}

	return conn.UpdateZone(setting)
}

func RemoveTrafficRule(zone string, rule *TrafficRule, permanent bool) error {
	if rule == nil {
		return nil
	}

	conn, err := getConn(permanent)
	if err != nil {
		return err
	}

	setting, err := conn.GetZoneByName(zone)
	if err != nil {
		return err
	}

	removeRichRule := false
	if rule.needRichRule() {
		removeRichRule = true
	} else {
		found := false
		switch rule.Type {
		case ServiceRuleType:
			var val string
			err := json.Unmarshal(rule.Value, &val)
			if err != nil {
				return err
			}

			setting.Services = slices.DeleteFunc(setting.Services, func(s string) bool {
				if s == val {
					found = true
					return true
				}

				return false
			})
		case PortRuleType:
			var port Port
			err := json.Unmarshal(rule.Value, &port)
			if err != nil {
				return err
			}

			setting.Ports = slices.DeleteFunc(setting.Ports, func(port2 *firewalld.Port) bool {
				if port.Port == port2.Port && port.Protocol == port2.Protocol {
					found = true
					return true
				}

				return false
			})
		case ProtocolRuleType:
			var val string
			err := json.Unmarshal(rule.Value, &val)
			if err != nil {
				return err
			}

			setting.Protocols = slices.DeleteFunc(setting.Protocols, func(s string) bool {
				if s == val {
					found = true
					return true
				}

				return false
			})
		case MasqueradeRuleType:
			if setting.Masquerade {
				setting.Masquerade = false
				found = true
			}
		case IcmpBlockRuleType:
			var val string
			err := json.Unmarshal(rule.Value, &val)
			if err != nil {
				return err
			}

			setting.ICMPBlocks = slices.DeleteFunc(setting.ICMPBlocks, func(s string) bool {
				if s == val {
					found = true
					return true
				}

				return false
			})
		case SourcePortRuleType:
			var port Port
			err := json.Unmarshal(rule.Value, &port)
			if err != nil {
				return err
			}

			setting.SourcePorts = slices.DeleteFunc(setting.SourcePorts, func(port2 *firewalld.Port) bool {
				if port.Port == port2.Port && port.Protocol == port2.Protocol {
					found = true
					return true
				}

				return false
			})
		}

		if !found {
			// not found, try to remove from richRules
			removeRichRule = true
		}
	}

	found := !removeRichRule
	if removeRichRule {
		richRule, err := rule.toRichRule()
		if err != nil {
			return err
		}
		setting.RichRules = slices.DeleteFunc(setting.RichRules, func(s string) bool {
			if s == richRule {
				found = true
				return true
			}

			return false
		})
	}

	if !found {
		return NotFoundErr
	}

	return conn.UpdateZone(setting)
}

func (t *TrafficRule) needRichRule() bool {
	if t.Strategy == AcceptRuleStrategy {
		return !(t.Family == "" && t.SrcAddr == "" && !t.SrcAddrInvert && t.DestAddr == "" && !t.DestAddrInvert && !t.Log.Enabled && !t.Audit)
	}

	return true
}

func (t *TrafficRule) toRichRule() (string, error) {
	richRule := "rule"
	if t.Family != "" {
		richRule += fmt.Sprintf(" family=\"%s\"", t.Family)
	}

	if t.SrcAddr != "" {
		if t.SrcAddrInvert {
			richRule += fmt.Sprintf(" source NOT address=\"%s\"", t.SrcAddr)
		} else {
			richRule += fmt.Sprintf(" source address=\"%s\"", t.SrcAddr)
		}
	}

	if t.DestAddr != "" {
		if t.DestAddrInvert {
			richRule += fmt.Sprintf(" destination NOT address=\"%s\"", t.DestAddr)
		} else {
			richRule += fmt.Sprintf(" destination address=\"%s\"", t.DestAddr)
		}
	}

	switch t.Type {
	case ServiceRuleType:
		richRule += " service name=" + string(t.Value)
	case PortRuleType:
		var port Port
		err := json.Unmarshal(t.Value, &port)
		if err != nil {
			return "", err
		}

		richRule += fmt.Sprintf(" port port=\"%s\" protocol=\"%s\"", port.Port, port.Protocol)
	case ProtocolRuleType:
		richRule += " protocol value=" + string(t.Value)

	case MasqueradeRuleType:
		richRule += " masquerade"
	case IcmpBlockRuleType:
		richRule += " icmp-block name=" + string(t.Value)
	case SourcePortRuleType:
		var port Port
		err := json.Unmarshal(t.Value, &port)
		if err != nil {
			return "", err
		}

		richRule += fmt.Sprintf(" source-port port=\"%s\" protocol=\"%s\"", port.Port, port.Protocol)
	}

	if t.Log.Enabled {
		richRule += " log"
		if t.Log.Prefix != "" {
			richRule += fmt.Sprintf(" prefix=\"%s\"", t.Log.Prefix)
		}
		if t.Log.Level != "" {
			richRule += fmt.Sprintf(" level=\"%s\"", t.Log.Level)
		}
		richRule += fmt.Sprintf(" limit value=\"%s\"", t.Log.Limit)
	}

	if t.Audit {
		richRule += " audit"
	}

	if t.Type != MasqueradeRuleType && t.Type != IcmpBlockRuleType {
		switch t.Strategy {
		case AcceptRuleStrategy:
			richRule += " accept"
		case RejectRuleStrategy:
			richRule += " reject"
		case DropRuleStrategy:
			richRule += " drop"
		}
	}

	return richRule, nil
}
