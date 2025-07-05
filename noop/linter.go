// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"dagger/noop/internal/dagger"
)

type GoLinter struct{}

func (*Noop) GoLinter() *GoLinter {
	return &GoLinter{}
}

func (*GoLinter) Lint(ctx context.Context, ctr *dagger.Container) *dagger.File {
	return dag.File("noop_lint_report.txt", "")
}
