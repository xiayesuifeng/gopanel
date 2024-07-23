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

func GetPortForwards(zone string, permanent bool) ([]*PortForward, error) {
	conn, err := getConn(permanent)
	if err != nil {
		return nil, err
	}

	forwards, err := conn.GetZoneForwardPorts(zone)
	if err != nil {
		return nil, err
	}

	result := make([]*PortForward, 0, len(forwards))
	for _, forward := range forwards {
		result = append(result, &PortForward{
			Port:      forward.Port,
			Protocol:  ForwardProtocol(forward.Protocol),
			ToPort:    forward.ToPort,
			ToAddress: forward.ToAddress,
		})
	}

	return result, nil
}