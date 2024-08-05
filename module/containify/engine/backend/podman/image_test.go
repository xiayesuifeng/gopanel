package podman

import (
	"context"
	"testing"
)

var imageName = "archlinux"

func TestImage_Pull(t *testing.T) {
	instance := &Podman{}
	err := instance.New(nil)
	if err != nil {
		t.Error(err)
	}

	imageID, err := instance.Image().Pull(context.TODO(), imageName, nil)
	if err != nil {
		t.Error(err)
	}

	t.Logf(imageID)
}

func TestImage_List(t *testing.T) {
	instance := &Podman{}
	err := instance.New(nil)
	if err != nil {
		t.Error(err)
	}

	list, err := instance.Image().List(context.TODO(), false, nil)
	if err != nil {
		t.Error(err)
	}

	for _, image := range list {
		t.Log(image)
	}
}

func TestImage_Exists(t *testing.T) {
	instance := &Podman{}
	err := instance.New(nil)
	if err != nil {
		t.Error(err)
	}

	exist, err := instance.Image().Exists(context.TODO(), imageName)
	if err != nil {
		t.Error(err)
	}

	t.Logf("image %s exist: %t", imageName, exist)
}

func TestImage_Remove(t *testing.T) {
	instance := &Podman{}
	err := instance.New(nil)
	if err != nil {
		t.Error(err)
	}

	err = instance.Image().Remove(context.TODO(), imageName)
	if err != nil {
		t.Error(err)
	}

	t.Logf("image %s remove success", imageName)
}
