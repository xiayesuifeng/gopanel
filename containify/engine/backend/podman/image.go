package podman

import (
	"context"
	"errors"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine/entity"
)

type image struct {
	podman *Podman
}

func (i *image) Remove(ctx context.Context, nameOrID string) error {
	conn, err := i.podman.getConn(ctx)
	if err != nil {
		return err
	}

	_, errs := images.Remove(conn, []string{nameOrID}, &images.RemoveOptions{})
	return errors.Join(errs...)
}

func (i *image) Exists(ctx context.Context, nameOrID string) (bool, error) {
	conn, err := i.podman.getConn(ctx)
	if err != nil {
		return false, err
	}

	return images.Exists(conn, nameOrID, &images.ExistsOptions{})
}

func (i *image) List(ctx context.Context, all bool, filters map[string][]string) ([]*entity.Image, error) {
	conn, err := i.podman.getConn(ctx)
	if err != nil {
		return nil, err
	}

	list, err := images.List(conn, &images.ListOptions{
		All:     &all,
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	var images []*entity.Image
	for _, image := range list {
		images = append(images, &entity.Image{
			ID:          image.ID,
			ParentID:    image.ParentId,
			RepoTags:    image.RepoTags,
			RepoDigests: image.RepoDigests,
			Created:     image.Created,
			Size:        image.Size,
			VirtualSize: image.VirtualSize,
			Labels:      image.Labels,
			Containers:  image.Containers,
		})
	}

	return images, nil
}
