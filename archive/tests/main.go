// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"dagger/archive-tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type ArchiveTests struct {
	// +private
	Archive *dagger.Archive
}

func New() *ArchiveTests {
	return &ArchiveTests{
		Archive: dag.Archive(),
	}
}

func (m *ArchiveTests) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(m.Zip().All)
	ep.Go(m.Tar().All)

	return ep.Wait()
}
