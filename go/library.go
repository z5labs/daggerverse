// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"errors"
	"fmt"

	"dagger/go/internal/dagger"
)

type Linter interface {
	DaggerObject

	Lint(ctx context.Context, ctr *dagger.Container) *dagger.File
}

type StaticAnalyzer interface {
	DaggerObject

	StaticAnalysis(
		ctx context.Context,
		ctr *dagger.Container,
		lintReport *dagger.File,
		coverageReport *dagger.File,
	) error
}

type Library struct {
	// +private
	Go *Go

	// +private
	Linter Linter

	// +private
	StaticAnalyzer StaticAnalyzer
}

// A set of functions for working with a library written in Go.
func (m *Go) Library(
	// The Go module source code for the library.
	module *dagger.Directory,

	// +optional
	linter Linter,

	// +optional
	staticAnalyzer StaticAnalyzer,
) *Library {
	return &Library{
		Go:             m.WithWorkdir("/src", module),
		Linter:         linter,
		StaticAnalyzer: staticAnalyzer,
	}
}

// Run all continuous integration functions.
func (lib *Library) Ci(ctx context.Context) error {
	err := lib.Generate(ctx, "./...")
	if err != nil {
		return err
	}

	err = lib.Tidy(ctx)
	if err != nil {
		return err
	}

	lintReport := lib.Lint(ctx)

	coverageReport := lib.Test("./...", true)

	err = lib.StaticAnalysis(ctx, lintReport, coverageReport)
	if err != nil {
		return err
	}

	return nil
}

// Run generate directives and validate no filesystem changes.
func (lib *Library) Generate(
	ctx context.Context,

	// +default="./..."
	pkg string,
) error {
	entries, err := lib.Go.Generate(pkg, nil).Diff(ctx)
	if err != nil {
		return err
	}

	if len(entries) > 0 {
		return fmt.Errorf("forgot to run go generate")
	}

	return nil
}

// Validate no necessary changes for go.mod or go.sum.
func (lib *Library) Tidy(ctx context.Context) error {
	diff, err := lib.Go.Tidy(nil).Diff(ctx)
	if err != nil {
		return err
	}

	if len(diff) != 0 {
		return errors.New("forgot to run go mod tidy")
	}

	return nil
}

// Lint source code.
func (lib *Library) Lint(ctx context.Context) *dagger.File {
	if lib.Linter == nil {
		return &dagger.File{}
	}

	return lib.Linter.Lint(ctx, lib.Go.Ctr)
}

// Run tests and return coverage report.
func (lib *Library) Test(
	// +default="./..."
	pkg string,

	// +default=true
	race bool,
) *dagger.File {
	return lib.Go.Test(pkg, nil, true).Coverage(Atomic)
}

// Perform static analysis.
func (lib *Library) StaticAnalysis(
	ctx context.Context,
	lintReport *dagger.File,
	coverageReport *dagger.File,
) error {
	if lib.StaticAnalyzer == nil {
		return nil
	}

	return lib.StaticAnalyzer.StaticAnalysis(
		ctx,
		lib.Go.Ctr,
		lintReport,
		coverageReport,
	)
}
