package main

import (
	"context"

	"dagger/tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Tests struct {
	// +private
	Go *dagger.Go
}

func New(
	// Defaults to latest for cgo support in tests.
	// +default="latest"
	version string,
) *Tests {
	return &Tests{
		dag.Go(dagger.GoOpts{Version: version}),
	}
}

func (m *Tests) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(m.Build().All)
	ep.Go(m.Test().All)
	ep.Go(m.Library().All)

	return ep.Wait()
}
