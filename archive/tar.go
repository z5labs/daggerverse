// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"dagger/archive/internal/archive"
	"dagger/archive/internal/dagger"
)

type Tar struct{}

func (m *Archive) Tar() *Tar {
	return &Tar{}
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

	_, err = file.Export(ctx, name)
	if err != nil {
		return nil, err
	}

	err = archive.ExtractTar(ctx, name, "out", gzip)
	if err != nil {
		return nil, err
	}

	return dag.CurrentModule().Workdir("out"), nil
}
