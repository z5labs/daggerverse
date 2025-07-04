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

// Tidy
type Tidy struct {
	// +private
	Ctr *dagger.Container
}

// Tidy
func (m *Go) Tidy(
	source *dagger.Directory,
) *Tidy {
	return &Tidy{
		Ctr: m.Ctr,
	}
}

func (t *Tidy) run(
	ctx context.Context,
	args ...string,
) (string, error) {
	cmd := []string{"go", "mod", "tidy"}
	cmd = append(cmd, args...)

	return t.Ctr.
		WithExec(cmd).
		Stdout(ctx)
}

// Report
func (t *Tidy) Report(ctx context.Context) ([]string, error) {
	stdout, err := t.run(ctx, "-v")
	if err != nil {
		return nil, err
	}

	modules := slices.Collect(strings.Lines(stdout))
	return modules, nil
}

// Diff
func (t *Tidy) Diff(ctx context.Context) ([]string, error) {
	stdout, err := t.run(ctx, "-diff")
	if err != nil {
		return nil, err
	}

	changes := slices.Collect(strings.Lines(stdout))
	return changes, nil
}
