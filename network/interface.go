package network

type Protocol int

const (
	ProtocolNone Protocol = iota
	ProtocolDHCP
	ProtocolStatic
	ProtocolPPPoE
)
