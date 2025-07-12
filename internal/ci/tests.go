// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"dagger/ci/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Tests struct{}

func (m *Ci) Tests() *Tests {
	return &Tests{}
}

func (t *Tests) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(dag.GoTests(dagger.GoTestsOpts{
		Version: "latest",
	}).All)

	ep.Go(dag.ArchiveTests().All)

	return ep.Wait()
}
