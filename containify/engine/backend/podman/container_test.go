package podman

import (
	"context"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine/entity"
	"testing"
)

func getInstance() (*Podman, error) {
	instance := &Podman{}
	err := instance.New(nil)
	return instance, err
}

func TestContainer_Create(t *testing.T) {
	instance, err := getInstance()
	if err != nil {
		t.Error(err)
	}

	if id, err := instance.Container().Create(context.TODO(), &entity.Container{
		ContainerBasic: entity.ContainerBasic{
			Image:   "archlinux",
			Name:    "test-create-c",
			Command: []string{"/usr/bin/bash", "-c", "sleep 600s"},
			Env: map[string]string{
				"TEST_ENV": "true",
			},
			Ports: []entity.PortMapping{
				{
					ContainerPort: 8080,
					HostPort:      8080,
				},
			},
		},
		Mounts: []entity.Mount{
			{
				Type:        entity.BindMountType,
				Destination: "/mnt",
				Source:      "/tmp",
				RW:          true,
			},
		},
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(id)
	}
}
