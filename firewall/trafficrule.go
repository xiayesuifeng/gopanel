package firewall

import (
	"encoding/json"
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
	Port string
	// Protocol tcp or udp
	Protocol string
}
