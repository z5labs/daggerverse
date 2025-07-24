// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"errors"
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
func (g *Generate) Diff(ctx context.Context) error {
	before := g.Ctr.Directory("/src")

	cmds, err := g.Report(ctx)
	if err != nil {
		return err
	}
	fmt.Println(strings.Join(cmds, "\n"))

	after := g.Ctr.Directory("/src")

	stdout, err := dag.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "diffutils"}).
		WithMountedDirectory("/before", before).
		WithMountedDirectory("/after", after).
		WithExec([]string{"diff", "-r", "-y", "--suppress-common-lines", "/before", "/after"}).
		Stdout(ctx)
	if err != nil {
		return err
	}

	fmt.Println(stdout)
	if stdout != "" {
		return errors.New("generate resulted in difference before and after")
	}

	return nil
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
