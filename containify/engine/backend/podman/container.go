package podman

import (
	"context"
	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine/entity"
)

type container struct {
	podman *Podman
}

func (c *container) Create(ctx context.Context, container *entity.Container) (containerID string, err error) {
	conn, err := c.podman.getConn(ctx)
	if err != nil {
		return
	}

	mounts := make([]specs.Mount, 0, len(container.Mounts))
	for _, mount := range container.Mounts {
		options := make([]string, 0)
		if !mount.RW {
			options = append(options, "ro")
		}

		mounts = append(mounts, specs.Mount{
			Destination: mount.Destination,
			Type:        string(mount.Type),
			Source:      mount.Source,
			Options:     options,
		})
	}

	portMappings := make([]nettypes.PortMapping, 0, len(container.Ports))
	for _, port := range container.Ports {
		portMappings = append(portMappings, nettypes.PortMapping{
			HostIP:        port.HostIP,
			ContainerPort: port.ContainerPort,
			HostPort:      port.HostPort,
			Range:         port.Range,
			Protocol:      port.Protocol,
		})
	}

	resp, err := containers.CreateWithSpec(conn, &specgen.SpecGenerator{
		ContainerBasicConfig: specgen.ContainerBasicConfig{
			Name:    container.Name,
			Command: container.Command,
			Env:     container.Env,
			Labels:  container.Labels,
		},
		ContainerStorageConfig: specgen.ContainerStorageConfig{
			Image:  container.Image,
			Mounts: mounts,
		},
		ContainerNetworkConfig: specgen.ContainerNetworkConfig{
			PortMappings: portMappings,
		},
	}, nil)

	containerID = resp.ID

	return
}

func (c *container) Remove(ctx context.Context, nameOrID string) error {
	conn, err := c.podman.getConn(ctx)
	if err != nil {
		return err
	}

	_, err = containers.Remove(conn, nameOrID, nil)
	return err
}

func (c *container) List(ctx context.Context) ([]*entity.ListContainer, error) {
	//TODO implement me
	panic("implement me")
}

func (c *container) Stop(ctx context.Context, nameOrID string) error {
	conn, err := c.podman.getConn(ctx)
	if err != nil {
		return err
	}

	return containers.Stop(conn, nameOrID, nil)
}

func (c *container) Start(ctx context.Context, nameOrID string) error {
	conn, err := c.podman.getConn(ctx)
	if err != nil {
		return err
	}

	return containers.Start(conn, nameOrID, nil)
}

func (c *container) Restart(ctx context.Context, nameOrID string) error {
	conn, err := c.podman.getConn(ctx)
	if err != nil {
		return err
	}

	return containers.Restart(conn, nameOrID, nil)
}
