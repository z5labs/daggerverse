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

// Library provides a set of functions for working with a library written in Go.
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

// Ci
func (l *Library) Ci(ctx context.Context) error {
	err := l.Generate(ctx, "./...")
	if err != nil {
		return err
	}

	err = l.Tidy(ctx)
	if err != nil {
		return err
	}

	lintReport := l.Lint(ctx)

	coverageReport := l.Test("./...", true)

	err = l.StaticAnalysis(ctx, lintReport, coverageReport)
	if err != nil {
		return err
	}

	return nil
}

// Generate
func (l *Library) Generate(
	ctx context.Context,

	// +default="./..."
	pkg string,
) error {
	entries, err := l.Go.Generate(pkg, nil).Diff(ctx)
	if err != nil {
		return err
	}
	fmt.Println("entries", entries)

	if len(entries) > 0 {
		return fmt.Errorf("forgot to run go generate")
	}

	return nil
}

// Tidy
func (l *Library) Tidy(ctx context.Context) error {
	diff, err := l.Go.Tidy(nil).Diff(ctx)
	if err != nil {
		return err
	}

	if len(diff) != 0 {
		return errors.New("forgot to run go mod tidy")
	}

	return nil
}

// Lint
func (l *Library) Lint(ctx context.Context) *dagger.File {
	if l.Linter == nil {
		return &dagger.File{}
	}

	return l.Linter.Lint(ctx, l.Go.Ctr)
}

// Test
func (l *Library) Test(
	// +default="./..."
	pkg string,

	// +default=true
	race bool,
) *dagger.File {
	return l.Go.Test(pkg, nil, true).Coverage(Atomic)
}

// StaticAnalysis
func (l *Library) StaticAnalysis(
	ctx context.Context,
	lintReport *dagger.File,
	coverageReport *dagger.File,
) error {
	if l.StaticAnalyzer == nil {
		return nil
	}

	return l.StaticAnalyzer.StaticAnalysis(
		ctx,
		l.Go.Ctr,
		lintReport,
		coverageReport,
	)
}
