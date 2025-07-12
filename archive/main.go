package main

import (
	"context"

	"dagger/archive/internal/dagger"

	"github.com/containerd/platforms"
)

type Archive struct {
	// +private
	Container *dagger.Container
}

func New(ctx context.Context) (*Archive, error) {
	ctrs, err := dag.Go().
		Module(dag.CurrentModule().Source()).
		Library().
		Application("./cmd/archive").
		Build(ctx, dagger.GoApplicationBuildOpts{
			Platforms: []dagger.Platform{dagger.Platform(platforms.DefaultString())},
		})
	if err != nil {
		return nil, err
	}

	return &Archive{
		Container: &ctrs[0],
	}, nil
}
