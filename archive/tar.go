// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"dagger/archive/internal/dagger"
)

type Tar struct {
	// +private
	Container *dagger.Container
}

func (m *Archive) Tar() *Tar {
	return &Tar{
		Container: m.Container,
	}
}

// Extract tar contents to a directory
func (t *Tar) Extract(
	ctx context.Context,

	file *dagger.File,

	// Enable gzip decompression.
	// +optional
	gzip bool,
) (*dagger.Directory, error) {
	name, err := file.Name(ctx)
	if err != nil {
		return nil, err
	}

	args := []string{"tar", "extract"}

	if gzip {
		args = append(args, "--gzip")
	}

	args = append(args, name, "out")

	dir := t.Container.
		WithFile(name, file).
		WithExec(args, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		Directory("out")

	return dir, nil
}
