// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"
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

	// The Go module source code.
	// +optional
	module *dagger.Directory,
) *Generate {
	ctr := m.Ctr
	if module != nil {
		ctr = ctr.WithMountedDirectory("/src", module).
			WithWorkdir("/src")
	}

	return &Generate{
		Ctr: ctr,
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

	g.Ctr = g.Ctr.WithExec(cmd)

	return g.Ctr.Stdout(ctx)
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

// Diff
func (g *Generate) Diff(ctx context.Context) ([]string, error) {
	before := g.Ctr.Directory("/src")

	cmds, err := g.Report(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println(strings.Join(cmds, "\n"))

	after := g.Ctr.Directory("/src")

	entries, err := before.Diff(after).Entries(ctx)
	if err != nil {
		fmt.Println("failed to diff:", err)
		return entries, nil
	}

	return entries, nil
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
