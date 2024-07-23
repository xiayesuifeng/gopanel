package firewall

type ForwardProtocol string

const (
	TCPForwardProtocol  ForwardProtocol = "tcp"
	UDPForwardProtocol  ForwardProtocol = "udp"
	SCTPForwardProtocol ForwardProtocol = "sctp"
	DCCPForwardProtocol ForwardProtocol = "dccp"
)

type PortForward struct {
	// Port port number or range
	Port     string          `json:"port"`
	Protocol ForwardProtocol `json:"protocol"`
	// ToPort port number or range
	ToPort    string `json:"toPort"`
	ToAddress string `json:"toAddress"`
}
