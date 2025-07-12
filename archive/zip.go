// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"dagger/archive/internal/dagger"
)

type Zip struct {
	// +private
	Container *dagger.Container
}

func (m *Archive) Zip() *Zip {
	return &Zip{
		Container: m.Container,
	}
}

// Extract zip contents to a directory
func (z *Zip) Extract(
	ctx context.Context,
	file *dagger.File,
) (*dagger.Directory, error) {
	name, err := file.Name(ctx)
	if err != nil {
		return nil, err
	}

	args := []string{"zip", "extract", name, "out"}

	dir := z.Container.
		WithFile(name, file).
		WithExec(args, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		Directory("out")

	return dir, nil
}
