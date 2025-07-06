// A generated module for Tests functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"

	"dagger/tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Tests struct{}

func (m *Tests) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(m.GoApplicationTest)

	return ep.Wait()
}

func (m *Tests) GoApplicationTest(ctx context.Context) error {
	app := dag.Go().Application(dag.CurrentModule().Source().Directory("testdata/goapp"), dagger.GoApplicationOpts{
		ContainerScanner: dag.Trivy().AsGoContainerScanner(),
	})

	variants, err := app.Build(ctx)
	if err != nil {
		return err
	}

	ctrs := make([]*dagger.Container, len(variants))
	for i := range variants {
		ctrs[i] = &variants[i]
	}

	return app.Scan(ctx, ctrs)
}
