package firewall

import (
	"encoding/json"
	"gitlab.com/xiayesuifeng/go-firewalld"
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
				if strings.HasPrefix(args[i+1], "prefix") {
					i++
					rule.Log.Prefix = strings.Trim(strings.Split(args[i], "=")[1], "\"")
				}
				if strings.HasPrefix(args[i+1], "level") {
					i++
					rule.Log.Level = strings.Trim(strings.Split(args[i], "=")[1], "\"")
				}
				if strings.HasPrefix(args[i+1], "limit") {
					i += 2
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

	needRichRule := rule.needRichRule()
	richRule := "rule"
	if needRichRule {
		if rule.Family != "" {
			richRule += " family=" + rule.Family
		}

		if rule.SrcAddr != "" {
			if rule.SrcAddrInvert {
				richRule += " source NOT address=" + rule.SrcAddr
			} else {
				richRule += " source address=" + rule.SrcAddr
			}
		}

		if rule.DestAddr != "" {
			if rule.DestAddrInvert {
				richRule += " destination NOT address=" + rule.DestAddr
			} else {
				richRule += " destination address=" + rule.DestAddr
			}
		}
	}

	switch rule.Type {
	case ServiceRuleType:
		if needRichRule {
			richRule += " service name=" + string(rule.Value)
		} else {
			var val string
			err := json.Unmarshal(rule.Value, &val)
			if err != nil {
				return err
			}
			setting.Services = append(setting.Services, val)
		}
	case PortRuleType:
		var port Port
		err := json.Unmarshal(rule.Value, &port)
		if err != nil {
			return err
		}

		if needRichRule {
			richRule += " port port=" + port.Port
			if port.Protocol != "" {
				richRule += " protocol=" + port.Protocol
			}
		} else {
			setting.Ports = append(setting.Ports, &firewalld.Port{
				Port:     port.Port,
				Protocol: port.Protocol,
			})
		}
	case ProtocolRuleType:
		if needRichRule {
			richRule += " protocol value=" + string(rule.Value)
		} else {
			var val string
			err := json.Unmarshal(rule.Value, &val)
			if err != nil {
				return err
			}

			setting.Protocols = append(setting.Protocols, val)
		}
	case MasqueradeRuleType:
		if needRichRule {
			richRule += " masquerade"
		} else {
			setting.Masquerade = true
		}
	case IcmpBlockRuleType:
		if needRichRule {
			richRule += " icmp-block name=" + string(rule.Value)
		} else {
			var val string
			err := json.Unmarshal(rule.Value, &val)
			if err != nil {
				return err
			}

			setting.ICMPBlocks = append(setting.Protocols, val)
		}
	case SourcePortRuleType:
		var port Port
		err := json.Unmarshal(rule.Value, &port)
		if err != nil {
			return err
		}

		if needRichRule {
			richRule += " source-port port=" + port.Port
			if port.Protocol != "" {
				richRule += " protocol=" + port.Protocol
			}
		} else {
			setting.SourcePorts = append(setting.SourcePorts, &firewalld.Port{
				Port:     port.Port,
				Protocol: port.Protocol,
			})
		}
	}

	if needRichRule {
		if rule.Log.Enabled {
			richRule += " log"
			if rule.Log.Prefix != "" {
				richRule += " prefix=" + rule.Log.Prefix
			}
			if rule.Log.Level != "" {
				richRule += " level=" + rule.Log.Level
			}
			if rule.Log.Prefix != "" {
				richRule += " prefix=" + rule.Log.Prefix
			}
			richRule += " limit value=" + rule.Log.Limit
		}

		if rule.Audit {
			richRule += " audit"
		}

		if rule.Type != MasqueradeRuleType && rule.Type != IcmpBlockRuleType {
			switch rule.Strategy {
			case AcceptRuleStrategy:
				richRule += " accept"
			case RejectRuleStrategy:
				richRule += " reject"
			case DropRuleStrategy:
				richRule += " drop"
			}
		}

		setting.RichRules = append(setting.RichRules, richRule)
	}

	return conn.UpdateZone(setting)
}

func (t *TrafficRule) needRichRule() bool {
	if t.Strategy == AcceptRuleStrategy {
		return !(t.Family == "" && t.SrcAddr == "" && !t.SrcAddrInvert && t.DestAddr == "" && !t.DestAddrInvert && !t.Log.Enabled && !t.Audit)
	}

	return true
}
