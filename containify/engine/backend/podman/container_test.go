package podman

import (
	"context"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine/entity"
	"testing"
)

const containerName = "test-create-c"

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
			Name:    containerName,
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

func TestContainer_List(t *testing.T) {
	instance, err := getInstance()
	if err != nil {
		t.Error(err)
	}

	list, err := instance.container.List(context.TODO())
	if err != nil {
		t.Error(err)
	}

	for _, container := range list {
		t.Log(container)
	}
}

func TestContainer_Start(t *testing.T) {
	instance, err := getInstance()
	if err != nil {
		t.Error(err)
	}

	err = instance.container.Start(context.TODO(), containerName)
	if err != nil {
		t.Error(err)
	}
}

func TestContainer_Restart(t *testing.T) {
	instance, err := getInstance()
	if err != nil {
		t.Error(err)
	}

	err = instance.container.Restart(context.TODO(), containerName)
	if err != nil {
		t.Error(err)
	}
}

func TestContainer_Stop(t *testing.T) {
	instance, err := getInstance()
	if err != nil {
		t.Error(err)
	}

	err = instance.container.Stop(context.TODO(), containerName)
	if err != nil {
		t.Error(err)
	}
}

func TestContainer_Remove(t *testing.T) {
	instance, err := getInstance()
	if err != nil {
		t.Error(err)
	}

	err = instance.container.Remove(context.TODO(), containerName)
	if err != nil {
		t.Error(err)
	}
}
