package firewall

import (
	"cmp"
	"errors"
	"gitlab.com/xiayesuifeng/go-firewalld"
	"slices"
	"sync"
)

type PolicyStrategy string

const (
	AcceptPolicyStrategy   = "ACCEPT"
	DefaultPolicyStrategy  = "default"
	RejectPolicyStrategy   = "REJECT"
	ContinuePolicyStrategy = "CONTINUE"
	DropPolicyStrategy     = "DROP"
)

type Policy struct {
	Name         string         `json:"name"`
	Short        string         `json:"short"`
	Description  string         `json:"description"`
	Target       string         `json:"target"`
	IngressZones []string       `json:"ingressZones"`
	EgressZones  []string       `json:"egressZones"`
	Services     []string       `json:"services"`
	ICMPBlocks   []string       `json:"icmpBlocks"`
	Priority     int            `json:"priority"`
	Masquerade   bool           `json:"masquerade"`
	ForwardPorts []*PortForward `json:"forwardPorts"`
	RichRules    []string       `json:"richRules"`
	Protocols    []string       `json:"protocols"`
	Ports        []*Port        `json:"ports"`
	SourcePorts  []*Port        `json:"sourcePorts"`
}

func GetPolicies(permanent bool) (result []*Policy, err error) {
	conn, err := getConn(permanent)
	if err != nil {
		return nil, err
	}

	names, err := conn.GetPolicyNames()
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
			policy, err := conn.GetPolicyByName(name)
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				errs = append(errs, err)
				return
			}

			result = append(result, toPolicy(policy))
		}()
	}

	wg.Wait()
	err = errors.Join(errs...)

	if err == nil {
		slices.SortFunc(result, func(a, b *Policy) int {
			return cmp.Compare[string](a.Name, b.Name)
		})
	}

	return
}

func toPolicy(policy *firewalld.Policy) *Policy {
	icmpBlocks := make([]string, 0)
	if policy.ICMPBlocks != nil {
		icmpBlocks = policy.ICMPBlocks
	}
	ingressZones := make([]string, 0)
	if policy.IngressZones != nil {
		ingressZones = policy.IngressZones
	}
	egressZones := make([]string, 0)
	if policy.EgressZones != nil {
		egressZones = policy.EgressZones
	}

	protocols := make([]string, 0)
	if policy.Protocols != nil {
		protocols = policy.Protocols
	}

	richRules := make([]string, 0)
	if policy.RichRules != nil {
		richRules = policy.RichRules
	}

	services := make([]string, 0)
	if policy.Services != nil {
		services = policy.Services
	}

	forwardPorts := make([]*PortForward, 0, len(policy.ForwardPorts))
	for _, port := range policy.ForwardPorts {
		forwardPorts = append(forwardPorts, &PortForward{
			Port:      port.Port,
			Protocol:  ForwardProtocol(port.Protocol),
			ToPort:    port.ToPort,
			ToAddress: port.ToAddress,
		})
	}

	ports := make([]*Port, 0, len(policy.Ports))
	for _, port := range policy.Ports {
		ports = append(ports, &Port{
			Port:     port.Port,
			Protocol: port.Protocol,
		})
	}

	sourcePorts := make([]*Port, 0, len(policy.SourcePorts))
	for _, port := range policy.SourcePorts {
		sourcePorts = append(sourcePorts, &Port{
			Port:     port.Port,
			Protocol: port.Protocol,
		})
	}

	return &Policy{
		Name:         policy.Name,
		Short:        policy.Short,
		Description:  policy.Description,
		Target:       policy.Target,
		IngressZones: ingressZones,
		EgressZones:  egressZones,
		Services:     protocols,
		ICMPBlocks:   icmpBlocks,
		Priority:     policy.Priority,
		Masquerade:   policy.Masquerade,
		ForwardPorts: forwardPorts,
		RichRules:    richRules,
		Protocols:    services,
		Ports:        ports,
		SourcePorts:  sourcePorts,
	}
}
