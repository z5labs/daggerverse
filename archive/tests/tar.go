// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"dagger/archive-tests/internal/dagger"
	"fmt"

	"github.com/sourcegraph/conc/pool"
)

type Tar struct {
	// +private
	Tar *dagger.ArchiveTar
}

func (m *ArchiveTests) Tar() *Tar {
	return &Tar{
		Tar: m.Archive.Tar(),
	}
}

func (t *Tar) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(t.FromUrlTest)

	return ep.Wait()
}

func (t *Tar) FromUrlTest(ctx context.Context) error {
	file := dag.HTTP("https://github.com/protocolbuffers/protobuf/releases/download/v31.1/protobuf-31.1.tar.gz")

	dir := t.Tar.Extract(file, dagger.ArchiveTarExtractOpts{
		Gzip: true,
	})

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 2 {
		return fmt.Errorf("expected 2 top level entries in directory instead of: %d", len(entries))
	}

	return nil
}
