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

type Zip struct{}

func (m *Archive) Zip() *Zip {
	return &Zip{}
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

	_, err = file.Export(ctx, name)
	if err != nil {
		return nil, err
	}

	err = archive.ExtractZip(ctx, name, "out")
	if err != nil {
		return nil, err
	}

	return dag.CurrentModule().Workdir("out"), nil
}
