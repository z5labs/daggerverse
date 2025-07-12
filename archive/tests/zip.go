// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"

	"dagger/archive-tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Zip struct {
	// +private
	Zip *dagger.ArchiveZip
}

func (m *ArchiveTests) Zip() *Zip {
	return &Zip{
		Zip: m.Archive.Zip(),
	}
}

func (z *Zip) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(z.FromUrlTest)

	return ep.Wait()
}

func (z *Zip) FromUrlTest(ctx context.Context) error {
	file := dag.HTTP("https://github.com/protocolbuffers/protobuf/releases/download/v31.1/protoc-31.1-linux-x86_64.zip")

	dir := z.Zip.Extract(file)

	entries, err := dir.Entries(ctx)
	if err != nil {
		return err
	}

	if len(entries) != 3 {
		return fmt.Errorf("expected 3 top level entries in directory instead of: %d", len(entries))
	}

	return nil
}
