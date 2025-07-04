// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"dagger/const/internal/dagger"
)

type GoLinter struct {
	// +private
	Report *dagger.File
}

func (*Const) GoLinter(report *dagger.File) *GoLinter {
	return &GoLinter{
		Report: report,
	}
}

func (g *GoLinter) GoLinter(ctx context.Context, ctr *dagger.Container) *dagger.File {
	return g.Report
}
