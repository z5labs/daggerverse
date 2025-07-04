// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"slices"
	"strings"

	"dagger/go/internal/dagger"
)

type Generate struct {
	// +private
	Ctr *dagger.Container

	// +private
	Pkg string
}

// Generate
func (m *Go) Generate(
	pkg string,
) *Generate {
	return &Generate{
		Ctr: m.Ctr,
		Pkg: pkg,
	}
}

func (g *Generate) run(
	ctx context.Context,
	pkg string,
	args ...string,
) (string, error) {
	cmd := []string{"go", "generate"}
	cmd = append(cmd, args...)
	cmd = append(cmd, pkg)

	return g.Ctr.
		WithExec(cmd).
		Stdout(ctx)
}

// Report
func (g *Generate) Report(ctx context.Context) ([]string, error) {
	stdout, err := g.run(ctx, g.Pkg, "-x")
	if err != nil {
		return nil, err
	}

	cmds := slices.Collect(strings.Lines(stdout))
	return cmds, nil
}

// DryRun
func (g *Generate) DryRun(ctx context.Context) ([]string, error) {
	stdout, err := g.run(ctx, g.Pkg, "-n")
	if err != nil {
		return nil, err
	}

	cmds := slices.Collect(strings.Lines(stdout))
	return cmds, nil
}
