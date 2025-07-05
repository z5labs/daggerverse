// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"errors"
	"fmt"

	"dagger/tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Build struct {
	// +private
	Go *dagger.Go
}

func (m *Tests) Build() *Build {
	return &Build{
		Go: m.Go,
	}
}

func (b *Build) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(b.WithDefaultsTest)
	ep.Go(b.WithOptionsTest)

	return ep.Wait()
}

func (b *Build) WithDefaultsTest(ctx context.Context) error {
	f := b.Go.Build(".", dagger.GoBuildOpts{
		Module: dag.CurrentModule().Source().Directory("testdata/buildoutput"),
	}).Output()

	n, err := f.Size(ctx)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("expected non-empty output file")
	}

	return nil
}

func (b *Build) WithOptionsTest(ctx context.Context) error {
	f := b.Go.Build(".", dagger.GoBuildOpts{
		Module: dag.CurrentModule().Source().Directory("testdata/buildoutput"),
	}).Output()

	withDebugSize, err := f.Size(ctx)
	if err != nil {
		return err
	}

	if withDebugSize == 0 {
		return errors.New("expected non-empty output file")
	}

	f = b.Go.Build(".", dagger.GoBuildOpts{
		Module:  dag.CurrentModule().Source().Directory("testdata/buildoutput"),
		Ldflags: []string{"-s", "-w"},
	}).Output()

	withoutDebugSize, err := f.Size(ctx)
	if err != nil {
		return err
	}

	if withoutDebugSize == 0 {
		return errors.New("expected non-empty output file")
	}

	if withDebugSize < withoutDebugSize {
		return fmt.Errorf(
			"expected output file with debug symbols to be larger than without: %d:%d",
			withDebugSize,
			withoutDebugSize,
		)
	}

	return nil
}
