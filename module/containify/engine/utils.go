package engine

import (
	"context"
	image2 "github.com/containers/image/v5/image"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"gitlab.com/xiayesuifeng/gopanel/module/containify/engine/entity"
)

func InspectImage(ctx context.Context, image string) (*entity.InspectImage, error) {
	sysCtx := &types.SystemContext{}

	ref, err := alltransports.ParseImageName("docker://" + image)
	if err != nil {
		return nil, err
	}

	src, err := ref.NewImageSource(ctx, sysCtx)
	if err != nil {
		return nil, err
	}

	img, err := image2.FromUnparsedImage(ctx, sysCtx, image2.UnparsedInstance(src, nil))
	if err != nil {
		return nil, err
	}

	config, err := img.OCIConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.InspectImage{
		ExposedPorts: config.Config.ExposedPorts,
		Env:          config.Config.Env,
		Entrypoint:   config.Config.Entrypoint,
		Cmd:          config.Config.Cmd,
		Volumes:      config.Config.Volumes,
		WorkingDir:   config.Config.WorkingDir,
		Labels:       config.Config.Labels,
	}, nil
}
