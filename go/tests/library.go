// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"errors"

	"dagger/gotests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Library struct {
	// +private
	Go *dagger.Go
}

func (m *GoTests) Library() *Library {
	return &Library{
		Go: m.Go,
	}
}

func (t *Library) All(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(t.CiTest)
	ep.Go(t.TidyTest)
	ep.Go(t.GenerateTest)
	ep.Go(t.GenerateContentDiffTest)
	ep.Go(t.GenerateNoDiffTest)

	return ep.Wait()
}

func (l *Library) CiTest(ctx context.Context) error {
	return l.Go.Module(dag.CurrentModule().Source().Directory("testdata/library/ci")).
		Library(
			dagger.GoModLibraryOpts{
				Linter:         dag.Noop().GoLinter().AsGoLinter(),
				StaticAnalyzer: dag.Noop().GoStaticAnalyzer().AsGoStaticAnalyzer(),
			},
		).
		Ci(ctx)
}

func (l *Library) TidyTest(ctx context.Context) error {
	err := l.Go.Module(dag.CurrentModule().Source().Directory("testdata/library/tidy")).
		Library(
			dagger.GoModLibraryOpts{
				Linter:         dag.Noop().GoLinter().AsGoLinter(),
				StaticAnalyzer: dag.Noop().GoStaticAnalyzer().AsGoStaticAnalyzer(),
			},
		).
		Tidy(ctx)

	if err == nil {
		return errors.New("expected tidy to fail due to missing deps")
	}

	return nil
}

func (l *Library) GenerateTest(ctx context.Context) error {
	err := l.Go.Module(dag.CurrentModule().Source().Directory("testdata/library/generate")).
		Library(
			dagger.GoModLibraryOpts{
				Linter:         dag.Noop().GoLinter().AsGoLinter(),
				StaticAnalyzer: dag.Noop().GoStaticAnalyzer().AsGoStaticAnalyzer(),
			},
		).
		Generate(ctx)

	if err == nil {
		return errors.New("expected generate to fail due to files being changed from go generate")
	}

	return nil
}

func (l *Library) GenerateContentDiffTest(ctx context.Context) error {
	err := l.Go.Module(dag.CurrentModule().Source().Directory("testdata/library/generate-content-diff")).
		Library(
			dagger.GoModLibraryOpts{
				Linter:         dag.Noop().GoLinter().AsGoLinter(),
				StaticAnalyzer: dag.Noop().GoStaticAnalyzer().AsGoStaticAnalyzer(),
			},
		).
		Generate(ctx)

	if err == nil {
		return errors.New("expected generate to fail due to files being changed from go generate")
	}

	return nil
}

func (l *Library) GenerateNoDiffTest(ctx context.Context) error {
	return l.Go.Module(dag.CurrentModule().Source().Directory("testdata/library/generate-no-diff")).
		Library(
			dagger.GoModLibraryOpts{
				Linter:         dag.Noop().GoLinter().AsGoLinter(),
				StaticAnalyzer: dag.Noop().GoStaticAnalyzer().AsGoStaticAnalyzer(),
			},
		).
		Generate(ctx)
}
