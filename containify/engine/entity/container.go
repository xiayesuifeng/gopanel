package entity

type ListContainer struct {
	ContainerBasic
}

type ContainerBasic struct {
	AutoRemove bool              `json:"autoRemove"`
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Image      string            `json:"image"`
	ImageID    string            `json:"imageID"`
	Command    []string          `json:"command"`
	Env        map[string]string `json:"env"`
	Ports      []PortMapping     `json:"ports"`
	Labels     map[string]string `json:"labels"`
	State      string            `json:"state"`
	Status     string            `json:"status"`
}

type Container struct {
	ContainerBasic
	Mounts []Mount
}

type MountType string

const (
	BindMountType   = "bind"
	VolumeMountType = "volume"
)

type Mount struct {
	// Whether the mount is a volume or bind mount.
	Type MountType `json:"type"`
	// The name of the volume. Empty for bind mounts.
	Name string `json:"Name,omitempty"`
	// The destination directory for the volume.
	Destination string `json:"destination"`
	// The driver used for the named volume. Empty for bind mounts.
	Driver string `json:"driver,omitempty"`
	// The source directory for the volume.
	Source string `json:"source"`
	// Whether the volume is read-write
	RW bool `json:"rw"`
}

type PortMapping struct {
	// HostIP is the IP that we will bind to on the host.
	// If unset, assumed to be 0.0.0.0 (all interfaces).
	HostIP string `json:"hostIP"`
	// ContainerPort is the port number that will be exposed from the
	// container.
	// Mandatory.
	ContainerPort uint16 `json:"containerPort"`
	// HostPort is the port number that will be forwarded from the host into
	// the container.
	// If omitted, a random port on the host (guaranteed to be over 1024)
	// will be assigned.
	HostPort uint16 `json:"hostPort"`
	// Range is the number of ports that will be forwarded, starting at
	// HostPort and ContainerPort and counting up.
	// This is 1-indexed, so 1 is assumed to be a single port (only the
	// Hostport:Containerport mapping will be added), 2 is two ports (both
	// Hostport:Containerport and Hostport+1:Containerport+1), etc.
	// If unset, assumed to be 1 (a single port).
	// Both hostport + range and containerport + range must be less than
	// 65536.
	Range uint16 `json:"range"`
	// Protocol is the protocol forward.
	// Must be either "tcp", "udp", and "sctp", or some combination of these
	// separated by commas.
	// If unset, assumed to be TCP.
	Protocol string `json:"protocol"`
}
