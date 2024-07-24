package firewall

import (
	"encoding/json"
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
