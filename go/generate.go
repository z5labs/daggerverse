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

// Run commands described by directives within existing files.
func (m *Mod) Generate(
	// Package path to search for directives within.
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

	g.Ctr = g.Ctr.WithExec(cmd)

	return g.Ctr.Stdout(ctx)
}

// Return all generate directives ran.
func (g *Generate) Report(ctx context.Context) ([]string, error) {
	stdout, err := g.run(ctx, g.Pkg, "-x")
	if err != nil {
		return nil, err
	}

	cmds := slices.Collect(strings.Lines(stdout))
	return cmds, nil
}

// Validate no change to the filesystem after running all generate directives.
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

// Return all generate directives which would be ran.
func (g *Generate) DryRun(ctx context.Context) ([]string, error) {
	stdout, err := g.run(ctx, g.Pkg, "-n")
	if err != nil {
		return nil, err
	}

	cmds := slices.Collect(strings.Lines(stdout))
	return cmds, nil
}
