package main

import (
	"context"

	"dagger/gotests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type GoTests struct {
	// +private
	Go *dagger.Go
}

func New(
	// Defaults to latest for cgo support in tests.
	// +default="latest"
	version string,
) *GoTests {
	return &GoTests{
		dag.Go(dagger.GoOpts{Version: version}),
	}
}

func (m *GoTests) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(m.Build().All)
	ep.Go(m.Test().All)
	ep.Go(m.Library().All)

	return ep.Wait()
}
